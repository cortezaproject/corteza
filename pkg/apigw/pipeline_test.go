package apigw

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func NewPl() *pl {
	return NewPipeline(zap.NewNop())
}

func Test_pipelineAdd(t *testing.T) {
	var (
		req = require.New(t)
	)

	p := NewPl()
	p.Add(mockExecer{})

	req.Len(p.w, 1)
}

func Test_pipelineExec(t *testing.T) {
	var (
		ctx   = context.Background()
		req   = require.New(t)
		scope = &scp{"foo": 1}
	)

	p := NewPl()
	p.Add(mockExecer{
		exec: func(c context.Context, s *scp) (err error) {
			s.Set("foo", 2)
			return nil
		},
	})

	err := p.Exec(ctx, scope)

	req.NoError(err)

	foo, err := scope.Get("foo")

	req.NoError(err)
	req.Equal(2, foo)
}

func Test_pipelineExecErr(t *testing.T) {
	var (
		ctx   = context.Background()
		req   = require.New(t)
		scope = &scp{"foo": 1}
	)

	p := NewPl()
	p.Add(mockExecer{
		exec: func(c context.Context, s *scp) (err error) {
			return fmt.Errorf("error returned")
		},
	})

	err := p.Exec(ctx, scope)

	req.Error(err, "error returned")
}
