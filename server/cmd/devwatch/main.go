// Command devwatch keeps the running dev server built from the code on disk.
//
// It watches Go sources under -root, rebuilds -build when they change, and
// replaces the running child. Two properties are the point of it, and both are
// load-bearing for anything that verifies a fix against this server:
//
//   - A change that lands DURING a build is not lost. Nothing here stamps a
//     clock against file mtimes; a write becomes a pending signal, and a signal
//     raised while a build is in flight is served by the cycle after it. The
//     tree is never re-read to decide what is new, so there is no window in
//     which a file can be too old to notice.
//
//   - Nothing ever runs a partly written binary. The build writes beside the
//     target and renames onto it, which is atomic within a directory, so the
//     child either starts from the previous build or from a complete new one.
//
// There is no proxy and no port of its own. Rebuilds are driven by the
// filesystem alone, so nothing has to make a request to provoke one.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// How long to gather further writes before building. A save from an editor is
// several events across several files, and each one is not its own build.
const settle = 300 * time.Millisecond

// How long the child gets to shut down on its own before it is killed.
const graceful = 5 * time.Second

func main() {
	var (
		root  = flag.String("root", ".", "directory tree to watch")
		build = flag.String("build", "./cmd/human", "package to build")
		bin   = flag.String("bin", "build/dev-bin", "where to write the binary, and what to run")
	)

	flag.Parse()

	logger := log.New(os.Stdout, "[devwatch] ", log.Ltime)

	release, err := claim(*bin)
	if err != nil {
		logger.Fatal(err)
	}
	defer release()

	if err := watch(*root, *build, *bin, flag.Args(), logger); err != nil {
		logger.Fatal(err)
	}
}

// claim takes an exclusive lock on the binary this watcher owns, so one
// checkout cannot end up with two watchers building over each other.
//
// A second one is otherwise silent and destructive: it renames its own build
// onto the same path, its child cannot bind the port the first child holds and
// exits, and both append to the same log. The lock is held by the process, so
// it is released by any death including a kill, and there is no stale file to
// clear.
func claim(bin string) (func(), error) {
	if err := os.MkdirAll(filepath.Dir(bin), 0o755); err != nil {
		return nil, fmt.Errorf("cannot prepare %s: %w", filepath.Dir(bin), err)
	}

	lock := bin + ".lock"

	f, err := os.OpenFile(lock, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", lock, err)
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()

		return nil, fmt.Errorf(
			"another watcher already has %s — stop it before starting a second, "+
				"or its build and yours will fight over the same binary", bin)
	}

	return func() { f.Close() }, nil
}

func watch(root, pkg, bin string, childArgs []string, logger *log.Logger) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("cannot watch the filesystem: %w", err)
	}
	defer watcher.Close()

	dirs, err := watchTree(watcher, root)
	if err != nil {
		return err
	}

	logger.Printf("watching %d directories under %s", dirs, root)

	// Buffered to one, and sent to without blocking: a write arriving while a
	// build runs leaves exactly one signal pending, and further writes collapse
	// into it. That pending signal is what makes a mid-build edit survive.
	changed := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}

				// A newly created package is watched from now on. Without this
				// a directory added after startup is invisible for the rest of
				// the session.
				if ev.Op&fsnotify.Create != 0 && isDir(ev.Name) {
					if _, err := watchTree(watcher, ev.Name); err != nil {
						logger.Printf("cannot watch %s: %v", ev.Name, err)
					}
				}

				if !isSource(ev.Name) {
					continue
				}

				// A permission change is not a change to the program, and
				// nothing else on this list is worth a 15s rebuild.
				if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
					continue
				}

				raise(changed)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Printf("watch error: %v", err)
			}
		}
	}()

	child := &server{bin: bin, args: childArgs, logger: logger}
	defer child.stop()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	stop := make(chan struct{})
	go func() {
		<-signals
		close(stop)
	}()

	rebuild(changed, stop, settle, func() { cycle(pkg, bin, child, logger) })

	logger.Print("stopping")

	return nil
}

// rebuild runs cycle, then runs it again for every change that has been raised
// since, until stop.
//
// The guarantee lives here. A change raised WHILE cycle is running is still
// pending when it returns, so the next pass serves it — the loop never decides
// what is new by re-reading the tree, which is what leaves a window a write can
// fall into.
func rebuild(changed <-chan struct{}, stop <-chan struct{}, settleFor time.Duration, cycle func()) {
	for {
		cycle()

		select {
		case <-stop:
			return
		case <-changed:
			// Anything written in the moment after the first event joins this
			// build rather than provoking another one.
			time.Sleep(settleFor)
			drain(changed)
		}
	}
}

// cycle builds, and on success swaps the running child for the new binary.
//
// A failed build leaves the previous child serving. That is a deliberate
// choice and not a comfortable one: it keeps a typo from taking the API down
// mid-edit, at the price of a server answering with code older than the tree.
// What makes it safe is that the staleness is detectable — the child's start
// time still predates the edit, which is what dev_server_status reports on.
func cycle(pkg, bin string, child *server, logger *log.Logger) {
	logger.Print("building")

	next := bin + ".next"

	out, err := exec.Command("go", "build", "-o", next, pkg).CombinedOutput()
	if err != nil {
		logger.Printf("BUILD FAILED, still serving the previous build: %v", err)
		if msg := strings.TrimSpace(string(out)); msg != "" {
			fmt.Fprintln(os.Stderr, msg)
		}
		return
	}

	child.stop()

	if err := os.Rename(next, bin); err != nil {
		logger.Printf("cannot install the new binary: %v", err)
		return
	}

	if err := child.start(); err != nil {
		logger.Printf("cannot start %s: %v", bin, err)
		return
	}

	logger.Printf("serving %s", bin)
}

// server is the child process being kept up to date.
type server struct {
	bin    string
	args   []string
	logger *log.Logger

	cmd  *exec.Cmd
	done chan struct{}
}

func (s *server) start() error {
	cmd := exec.Command(s.bin, s.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = childAttr()

	if err := cmd.Start(); err != nil {
		return err
	}

	// Waited on immediately, so a server that dies on its own is reaped rather
	// than left as a zombie until the next rebuild. Its exit is not logged
	// here: every stop() produces one, and the server's own output already
	// says what happened when the exit was not asked for.
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	s.cmd, s.done = cmd, done

	return nil
}

// stop asks the child to shut down, and insists if it will not.
func (s *server) stop() {
	if s.cmd == nil || s.cmd.Process == nil {
		return
	}

	pgid := -s.cmd.Process.Pid

	_ = syscall.Kill(pgid, syscall.SIGINT)

	select {
	case <-s.done:
	case <-time.After(graceful):
		s.logger.Printf("child did not stop in %s, killing it", graceful)
		_ = syscall.Kill(pgid, syscall.SIGKILL)
		<-s.done
	}

	s.cmd, s.done = nil, nil
}

// watchTree adds every directory under root that can hold source worth
// rebuilding for, and reports how many it added.
func watchTree(watcher *fsnotify.Watcher, root string) (int, error) {
	added := 0

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		// The root is watched whatever it is called. Testing it like any other
		// directory makes a relative root of "." skip the entire tree, and a
		// watcher that watches nothing looks exactly like one with nothing to
		// do.
		if path != root && skipDir(d.Name()) {
			return filepath.SkipDir
		}

		if err := watcher.Add(path); err != nil {
			return fmt.Errorf("cannot watch %s: %w", path, err)
		}

		added++

		return nil
	})

	return added, err
}

// skipDir names the directories a rebuild can never depend on.
//
// build is here for a reason beyond noise: the binary this writes lands in it,
// so watching it would make every successful build trigger the next one.
func skipDir(name string) bool {
	switch name {
	case "vendor", "build", "node_modules", ".git", "testdata", "var":
		return true
	}

	return strings.HasPrefix(name, ".")
}

// isSource decides what is worth a rebuild. Editors write their own scratch
// files next to yours — a swapfile is not a change to the program.
func isSource(path string) bool {
	base := filepath.Base(path)

	if strings.HasPrefix(base, ".") || strings.HasSuffix(base, "~") {
		return false
	}

	return strings.HasSuffix(base, ".go")
}

func isDir(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

// raise records that something changed, without ever blocking the watcher.
//
// A full channel already says "rebuild when you can", which is the same thing
// this call would say, so dropping the send loses nothing.
func raise(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func drain(ch <-chan struct{}) {
	select {
	case <-ch:
	default:
	}
}
