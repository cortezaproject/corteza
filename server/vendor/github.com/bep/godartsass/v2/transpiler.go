// Package godartsass provides a Go API for the Dass Sass Embedded protocol.
//
// Use the Start function to create and start a new thread safe transpiler.
// Close it when done.
package godartsass

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/cli/safeexec"

	"github.com/bep/godartsass/v2/internal/embeddedsass"
	"google.golang.org/protobuf/proto"
)

const defaultDartSassBinaryFilename = "sass"

// ErrShutdown will be returned from Execute and Close if the transpiler is or
// is about to be shut down.
var ErrShutdown = errors.New("connection is shut down")

// Start creates and starts a new SCSS transpiler that communicates with the
// Dass Sass Embedded protocol via Stdin and Stdout.
//
// Closing the transpiler will shut down the process.
//
// Note that the Transpiler is thread safe, and the recommended way of using
// this is to create one and use that for all the SCSS processing needed.
func Start(opts Options) (*Transpiler, error) {
	if err := opts.init(); err != nil {
		return nil, err
	}

	// See https://github.com/golang/go/issues/38736
	bin, err := safeexec.LookPath(opts.DartSassEmbeddedFilename)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(bin)
	cmd.Args = append(cmd.Args, "--embedded")
	cmd.Stderr = os.Stderr

	conn, err := newConn(cmd)
	if err != nil {
		return nil, err
	}

	if err := conn.Start(); err != nil {
		return nil, err
	}

	t := &Transpiler{
		opts:    opts,
		conn:    conn,
		lenBuf:  make([]byte, binary.MaxVarintLen64),
		idBuf:   make([]byte, binary.MaxVarintLen64),
		pending: make(map[uint32]*call),
	}

	go t.input()

	return t, nil
}

// Version returns version information about the Dart Sass frameworks used
// in dartSassEmbeddedFilename.
func Version(dartSassEmbeddedFilename string) (DartSassVersion, error) {
	var v DartSassVersion
	bin, err := safeexec.LookPath(dartSassEmbeddedFilename)
	if err != nil {
		return v, err
	}

	cmd := exec.Command(bin, "--embedded", "--version")
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return v, err
	}

	if err := json.Unmarshal(out, &v); err != nil {
		return v, err
	}

	return v, nil
}

type DartSassVersion struct {
	ProtocolVersion       string `json:"protocolVersion"`
	CompilerVersion       string `json:"compilerVersion"`
	ImplementationVersion string `json:"implementationVersion"`
	ImplementationName    string `json:"implementationName"`
	ID                    int    `json:"id"`
}

// Transpiler controls transpiling of SCSS into CSS.
type Transpiler struct {
	opts Options

	// stdin/stdout of the Dart Sass protocol
	conn   byteReadWriteCloser
	lenBuf []byte
	idBuf  []byte
	msgBuf []byte

	closing  bool
	shutdown bool

	// Protects the sending of messages to Dart Sass.
	sendMu sync.Mutex

	mu      sync.Mutex // Protects all below.
	seq     uint32
	pending map[uint32]*call
}

// IsShutDown checks if all pending calls have been shut down.
// Used in tests.
func (t *Transpiler) IsShutDown() bool {
	for _, p := range t.pending {
		if p.Error != ErrShutdown {
			return false
		}
	}
	return true
}

// Result holds the result returned from Execute.
type Result struct {
	CSS       string
	SourceMap string
}

// SassError is the error returned from Execute on compile errors.
type SassError struct {
	Message string `json:"message"`
	Span    struct {
		Text  string `json:"text"`
		Start struct {
			Offset int `json:"offset"`
			Column int `json:"column"`
		} `json:"start"`
		End struct {
			Offset int `json:"offset"`
			Column int `json:"column"`
		} `json:"end"`
		Url     string `json:"url"`
		Context string `json:"context"`
	} `json:"span"`
}

func (e SassError) Error() string {
	span := e.Span
	file := path.Clean(strings.TrimPrefix(span.Url, "file:"))
	return fmt.Sprintf("file: %q, context: %q: %s", file, span.Context, e.Message)
}

// Close closes the stream to the embedded Dart Sass Protocol, shutting it down.
// If it is already shutting down, ErrShutdown is returned.
func (t *Transpiler) Close() error {
	t.sendMu.Lock()
	defer t.sendMu.Unlock()
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closing {
		return ErrShutdown
	}

	t.closing = true
	err := t.conn.Close()

	if eerr, ok := err.(*exec.ExitError); ok {
		if eerr.ExitCode() == 1 {
			// This is the expected exit code when shutting down.
			return ErrShutdown
		}
	}

	return err
}

// Execute transpiles the string Source given in Args into CSS.
// If Dart Sass resturns a "compile failure", the error returned will be
// of type SassError.
func (t *Transpiler) Execute(args Args) (Result, error) {
	var result Result

	createInboundMessage := func(seq uint32) (*embeddedsass.InboundMessage, error) {
		if err := args.init(seq, t.opts); err != nil {
			return nil, err
		}

		message := &embeddedsass.InboundMessage_CompileRequest_{
			CompileRequest: &embeddedsass.InboundMessage_CompileRequest{
				Importers: args.sassImporters,
				Style:     args.sassOutputStyle,
				Input: &embeddedsass.InboundMessage_CompileRequest_String_{
					String_: &embeddedsass.InboundMessage_CompileRequest_StringInput{
						Syntax: args.sassSourceSyntax,
						Source: args.Source,
						Url:    args.URL,
					},
				},
				SourceMap:               args.EnableSourceMap,
				SourceMapIncludeSources: args.SourceMapIncludeSources,
			},
		}

		return &embeddedsass.InboundMessage{
			Message: message,
		}, nil
	}

	call, err := t.newCall(createInboundMessage, args)
	if err != nil {
		return result, err
	}

	select {
	case call = <-call.Done:
	case <-time.After(t.opts.Timeout):
		return result, errors.New("timeout waiting for Dart Sass to respond; note that this project is only compatible with the Dart Sass Binary found here: https://github.com/sass/dart-sass/releases/")
	}

	if call.Error != nil {
		return result, call.Error
	}

	response := call.Response
	csp := response.Message.(*embeddedsass.OutboundMessage_CompileResponse_)

	switch resp := csp.CompileResponse.Result.(type) {
	case *embeddedsass.OutboundMessage_CompileResponse_Success:
		result.CSS = resp.Success.Css
		result.SourceMap = resp.Success.SourceMap
	case *embeddedsass.OutboundMessage_CompileResponse_Failure:
		asJson, err := json.Marshal(resp.Failure)
		if err != nil {
			return result, err
		}
		var sassErr SassError
		err = json.Unmarshal(asJson, &sassErr)
		if err != nil {
			return result, err
		}
		return result, sassErr
	default:
		return result, fmt.Errorf("unsupported response type: %T", resp)
	}

	return result, nil
}

func (t *Transpiler) getCall(id uint32) *call {
	t.mu.Lock()
	defer t.mu.Unlock()
	call, found := t.pending[id]
	if !found {
		panic(fmt.Sprintf("call with ID %d not found", id))
	}
	return call
}

func (t *Transpiler) input() {
	var err error

	for err == nil {
		// The header is the length in bytes of the remaining message including the compilation ID.
		var l uint64

		l, err = binary.ReadUvarint(t.conn)
		if err != nil {
			break
		}

		plen := int(l)
		if len(t.msgBuf) < plen {
			t.msgBuf = make([]byte, plen)
		}

		buf := t.msgBuf[:plen]

		_, err = io.ReadFull(t.conn, buf)
		if err != nil {
			break
		}

		v, n := binary.Uvarint(buf)
		if n <= 0 {
			break
		}
		compilationID := uint32(v)

		buf = buf[n:]

		var msg embeddedsass.OutboundMessage

		if err = proto.Unmarshal(buf, &msg); err != nil {
			break
		}

		switch c := msg.Message.(type) {
		case *embeddedsass.OutboundMessage_CompileResponse_:
			// Attach it to the correct pending call.
			t.mu.Lock()
			call := t.pending[compilationID]
			delete(t.pending, compilationID)
			t.mu.Unlock()
			if call == nil {
				err = fmt.Errorf("call with ID %d not found", compilationID)
				break
			}
			call.Response = &msg
			call.done()
		case *embeddedsass.OutboundMessage_CanonicalizeRequest_:
			call := t.getCall(compilationID)
			resolved, resolveErr := call.importResolver.CanonicalizeURL(c.CanonicalizeRequest.GetUrl())

			var response *embeddedsass.InboundMessage_CanonicalizeResponse
			if resolveErr != nil {
				response = &embeddedsass.InboundMessage_CanonicalizeResponse{
					Id: c.CanonicalizeRequest.GetId(),
					Result: &embeddedsass.InboundMessage_CanonicalizeResponse_Error{
						Error: resolveErr.Error(),
					},
				}
			} else {
				var url *embeddedsass.InboundMessage_CanonicalizeResponse_Url
				if resolved != "" {
					url = &embeddedsass.InboundMessage_CanonicalizeResponse_Url{
						Url: resolved,
					}
				}
				response = &embeddedsass.InboundMessage_CanonicalizeResponse{
					Id:     c.CanonicalizeRequest.GetId(),
					Result: url,
				}
			}

			err = t.sendInboundMessage(
				compilationID,
				&embeddedsass.InboundMessage{
					Message: &embeddedsass.InboundMessage_CanonicalizeResponse_{
						CanonicalizeResponse: response,
					},
				},
			)
		case *embeddedsass.OutboundMessage_ImportRequest_:
			call := t.getCall(compilationID)
			url := c.ImportRequest.GetUrl()
			imp, loadErr := call.importResolver.Load(url)
			sourceSyntax := embeddedsass.Syntax_value[string(imp.SourceSyntax)]

			var response *embeddedsass.InboundMessage_ImportResponse
			var sourceMapURL string

			// Dart Sass expect a browser-accessible URL or an empty string.
			// If no URL is supplied, a `data:` URL wil be generated
			// automatically from `contents`
			if hasScheme(url) {
				sourceMapURL = url
			}

			if loadErr != nil {
				response = &embeddedsass.InboundMessage_ImportResponse{
					Id: c.ImportRequest.GetId(),
					Result: &embeddedsass.InboundMessage_ImportResponse_Error{
						Error: loadErr.Error(),
					},
				}
			} else {
				response = &embeddedsass.InboundMessage_ImportResponse{
					Id: c.ImportRequest.GetId(),
					Result: &embeddedsass.InboundMessage_ImportResponse_Success{
						Success: &embeddedsass.InboundMessage_ImportResponse_ImportSuccess{
							Contents:     imp.Content,
							SourceMapUrl: &sourceMapURL,
							Syntax:       embeddedsass.Syntax(sourceSyntax),
						},
					},
				}
			}

			err = t.sendInboundMessage(
				compilationID,
				&embeddedsass.InboundMessage{
					Message: &embeddedsass.InboundMessage_ImportResponse_{
						ImportResponse: response,
					},
				},
			)
		case *embeddedsass.OutboundMessage_LogEvent_:
			if t.opts.LogEventHandler != nil {
				var logEvent LogEvent
				e := c.LogEvent
				if e.Span != nil {
					u := e.Span.Url
					if u == "" {
						u = "stdin"
					}
					u, _ = url.QueryUnescape(u)
					logEvent = LogEvent{
						Type:    LogEventType(e.Type),
						Message: fmt.Sprintf("%s:%d:%d: %s", u, e.Span.Start.Line, e.Span.Start.Column, c.LogEvent.GetMessage()),
					}
				} else {
					logEvent = LogEvent{
						Type:    LogEventType(e.Type),
						Message: e.GetMessage(),
					}
				}

				t.opts.LogEventHandler(logEvent)

			}

		case *embeddedsass.OutboundMessage_Error:
			err = fmt.Errorf("SASS error: %s", c.Error.GetMessage())
		default:
			err = fmt.Errorf("unsupported response message type. %T", msg.Message)
		}

	}

	// Terminate pending calls.
	t.sendMu.Lock()
	defer t.sendMu.Unlock()
	t.mu.Lock()
	defer t.mu.Unlock()

	t.shutdown = true
	isEOF := err == io.EOF || strings.Contains(err.Error(), "already closed")
	if isEOF {
		if t.closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}

	for _, call := range t.pending {
		call.Error = err
		call.done()
	}
}

func (t *Transpiler) nextSeq() uint32 {
	t.seq++
	// The compilation ID 0 is reserved for `VersionRequest` and `VersionResponse`,
	// 4294967295 is reserved for error handling. This is the maximum number representable by a `uint32` so it should be safe to start over.
	if t.seq == 0 || t.seq == 4294967295 {
		t.seq = 1
	}
	return t.seq
}

func (t *Transpiler) newCall(createInbound func(seq uint32) (*embeddedsass.InboundMessage, error), args Args) (*call, error) {
	t.mu.Lock()
	id := t.nextSeq()
	req, err := createInbound(id)
	if err != nil {
		t.mu.Unlock()
		return nil, err
	}

	call := &call{
		Request:        req,
		Done:           make(chan *call, 1),
		importResolver: args.ImportResolver,
	}

	if t.shutdown || t.closing {
		t.mu.Unlock()
		call.Error = ErrShutdown
		call.done()
		return call, nil
	}

	t.pending[id] = call

	t.mu.Unlock()

	switch call.Request.Message.(type) {
	case *embeddedsass.InboundMessage_CompileRequest_:
	default:
		return nil, fmt.Errorf("unsupported request message type. %T", call.Request.Message)
	}

	return call, t.sendInboundMessage(id, call.Request)
}

func (t *Transpiler) sendInboundMessage(compilationID uint32, message *embeddedsass.InboundMessage) error {
	t.sendMu.Lock()
	defer t.sendMu.Unlock()
	t.mu.Lock()
	if t.closing || t.shutdown {
		t.mu.Unlock()
		return ErrShutdown
	}
	t.mu.Unlock()

	out, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %s", err)
	}

	// Every message must begin with a varint indicating the length in bytes of
	// the remaining message including the compilation ID
	reqLen := uint64(len(out))
	compilationIDLen := binary.PutUvarint(t.idBuf, uint64(compilationID))
	headerLen := binary.PutUvarint(t.lenBuf, reqLen+uint64(compilationIDLen))
	_, err = t.conn.Write(t.lenBuf[:headerLen])
	if err != nil {
		return err
	}
	_, err = t.conn.Write(t.idBuf[:compilationIDLen])
	if err != nil {
		return err
	}

	headerLen, err = t.conn.Write(out)
	if headerLen != len(out) {
		return errors.New("failed to write payload")
	}

	return err
}

type call struct {
	Request        *embeddedsass.InboundMessage
	Response       *embeddedsass.OutboundMessage
	importResolver ImportResolver

	Error error
	Done  chan *call
}

func (call *call) done() {
	select {
	case call.Done <- call:
	default:
	}
}

func hasScheme(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Scheme != ""
}
