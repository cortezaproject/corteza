# third_party

Dependencies this repo forks and builds from source, wired in through a
`replace` in `go.mod`. `go mod vendor` regenerates `vendor/` from here, so a
patch survives a re-vendor.

## gval

`github.com/PaesslerAG/gval`, forked from `github.com/cortezaproject/gval`
v1.2.4. One behavioural difference from upstream: a name registered with
`gval.Function` is the function only where the expression calls it —
`split(a, b)` is the function, a bare `split` is a variable. Upstream claims
the name either way, which puts every function name out of reach as a module
field, workflow variable or DAL attribute.
