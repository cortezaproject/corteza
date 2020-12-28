package wfexec

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	gwTestStep struct{}
)

func (*gwTestStep) Exec(context.Context, *ExecRequest) (ExecResponse, error) {
	return nil, nil
}

func TestGatewayPath(t *testing.T) {
	var (
		req = require.New(t)
		gwp *GatewayPath
		err error
	)

	gwp, err = NewGatewayPath(nil, "")
	req.NoError(err)
	req.NotNil(gwp)

	gwp, err = NewGatewayPath(nil, "a > 1")
	req.NoError(err)
	req.NotNil(gwp)

	gwp, err = NewGatewayPath(nil, "<>")
	req.Error(err)
}

func TestJoinGateway(t *testing.T) {
	var (
		req        = require.New(t)
		p1, p2, p3 = &wfTestStep{name: "p1"}, &wfTestStep{name: "p2"}, &wfTestStep{name: "p3"}
		gw         = JoinGateway(p1, p2, p3)

		r   ExecResponse
		err error
	)

	r, err = gw.Exec(nil, &ExecRequest{Caller: p1})
	req.NoError(err)
	req.Equal(&Joined{}, r)

	r, err = gw.Exec(nil, &ExecRequest{Caller: p2})
	req.NoError(err)
	req.Equal(&Joined{}, r)

	r, err = gw.Exec(nil, &ExecRequest{Caller: p3})
	req.NoError(err)
	req.IsType(Variables{}, r)
}

func TestForkGateway(t *testing.T) {
	var (
		req = require.New(t)
		gw  = ForkGateway()
	)

	r, err := gw.Exec(nil, nil)
	req.NoError(err)
	req.Equal(Steps{}, r)
	req.Empty(r)
}

func TestInclGateway(t *testing.T) {
	var (
		req = require.New(t)

		s1, s2, s3 = &wfTestStep{name: "s1"}, &wfTestStep{name: "s2"}, &wfTestStep{name: "s3"}
		gwp1, _    = NewGatewayPath(s1, "a > 10")
		gwp2, _    = NewGatewayPath(s2, "a > 5")
		gwp3, _    = NewGatewayPath(s3, "a > 0")

		gw, err = InclGateway(gwp1, gwp2, gwp3)
	)

	r, err := gw.Exec(context.Background(), &ExecRequest{Scope: Variables{"a": 11}})
	req.NoError(err)
	req.Equal(Steps{s1, s2, s3}, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: Variables{"a": 6}})
	req.NoError(err)
	req.Equal(Steps{s2, s3}, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: Variables{"a": 1}})
	req.NoError(err)
	req.Equal(Steps{s3}, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: Variables{"a": 0}})
	req.Error(err)
	req.Nil(r)
}

func TestExclGateway(t *testing.T) {
	var (
		req = require.New(t)

		s1, s2, s3 = &wfTestStep{name: "s1"}, &wfTestStep{name: "s2"}, &wfTestStep{name: "s3"}
		gwp1, _    = NewGatewayPath(s1, "a > 10")
		gwp2, _    = NewGatewayPath(s2, "a > 5")
		gwp3, _    = NewGatewayPath(s3, "")

		gw, err = ExclGateway(gwp1, gwp2, gwp3)
	)

	r, err := gw.Exec(context.Background(), &ExecRequest{Scope: Variables{"a": 11}})
	req.NoError(err)
	req.Equal(s1, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: Variables{"a": 6}})
	req.NoError(err)
	req.Equal(s2, r)

	r, err = gw.Exec(context.Background(), &ExecRequest{Scope: Variables{"a": 1}})
	req.NoError(err)
	req.Equal(s3, r)
}
