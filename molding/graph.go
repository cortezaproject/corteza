package molding

import (
	"context"
	"fmt"
	"sync"
)

type (
	Nodes []Node

	Node interface {
		NodeRef() string
	}

	Iterator interface {
		Node

		// Returns next node
		Next() Node

		// Sets next node
		//
		// When configuring complex workflows with
		// loops & join gateways we need explicitly
		// add edges
		SetNext(Node)
	}

	Joiner interface {
		Iterator

		// Consumes parent node and scope variables
		//
		// Returns non-nil when Join fn is called with all parents
		Join(Node, Variables) (Node, Variables, error)

		// Returns all possible input paths
		// @todo think about naming here,
		//       we have Next() that returns one node and Paths() that returns multiple
		Paths() Nodes
	}

	Tester interface {
		Node

		// Returns one or more nodes that satisfy the configured conditions
		//
		// Always returns at least one node
		Test(context.Context, Variables) (Nodes, error)

		// Returns all possible output paths
		// @todo think about naming here,
		//       we have Next() that returns one node and Paths() that returns multiple
		Paths() Nodes
	}

	Executor interface {
		Iterator
		Exec(context.Context, Variables) (Variables, error)
	}

	Finalizer interface {
		Node
		Finalize(context.Context, Variables) error
	}

	queue chan *payload

	payload struct {
		node  Node
		prev  Node
		scope Variables
	}

	workflow struct {
		wg sync.WaitGroup
		l  sync.Mutex

		queue queue
		err   chan error

		counter map[Node]int
	}
)

func (q queue) push(s Variables, p Node, nn ...Node) {
	for _, n := range nn {
		q <- &payload{scope: s, prev: p, node: n}
	}
}

func Workflow(ctx context.Context, start Node, scope Variables) error {
	if err := CheckNodes(start); err != nil {
		return err
	}

	var (
		e = &workflow{
			queue:   make(queue, 512),
			err:     make(chan error, 1),
			counter: make(map[Node]int),
		}
	)

	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.loop(ctx)
	}()

	e.queue.push(scope, nil, start)
	e.wg.Wait()

	select {
	case err := <-e.err:
		return err
	default:
		return nil
	}
}

func (e *workflow) loop(ctx context.Context) {
	for {
		select {
		case err := <-e.err:
			e.err <- err
			return

		case p := <-e.queue:
			if p == nil {
				// processing queue closed
				return
			}

			e.wg.Add(1)
			func() {
				defer e.wg.Done()
				if err := e.proc(ctx, p.prev, p.node, p.scope); err != nil {
					e.err <- err
				}
			}()

		}
	}
}

func (e *workflow) proc(ctx context.Context, p Node, n Node, scope Variables) error {
	e.l.Lock()
	defer e.l.Unlock()
	e.counter[n]++

	switch c := n.(type) {
	case Joiner:
		next, joinedScope, err := c.Join(p, scope)
		if err != nil {
			return err
		}

		if next == nil {
			// no next node, waiting for all paths to execute
			return nil
		}

		e.queue.push(joinedScope, c, c.Next())

	case Finalizer:
		if err := c.Finalize(ctx, scope); err != nil {
			return err
		}

		close(e.queue)

	case Tester:
		paths, err := c.Test(ctx, scope)
		if err != nil {
			return err
		}
		e.queue.push(scope, c, paths...)

	case Executor:
		results, err := c.Exec(ctx, scope)
		if err != nil {
			return err
		}
		e.queue.push(scope.Merge(results), c, c.Next())

	default:
		return fmt.Errorf("useless node")
	}

	return nil
}

func (nn Nodes) String() string {
	path := ""
	for i := range nn {
		if i > 0 {
			path += "."
		}
		path += nn[i].NodeRef()
	}

	return path
}
