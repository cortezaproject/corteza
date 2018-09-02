# ctx

Composable utilities for Go contexts.

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/SentimensRG/ctx)
[![Go Report Card](https://goreportcard.com/badge/github.com/SentimensRG/ctx?style=flat-square)](https://goreportcard.com/report/github.com/SentimensRG/ctx)

## Installation

```bash
go get -u github.com/SentimensRG/ctx
```

## Overview

The `ctx` package provides utilites for working with data structures satisfying
the `ctx.Doner` interface, most notably `context.Context`:

```go
type Doner interface {
    Done() <-chan struct{}
}
```

The functions in `ctx` are appropriate for operations that do not preserve the
values in a context, e.g.: joining several contexts together.

## Subpackages

- [sigctx](https://github.com/SentimensRG/ctx/tree/master/sigctx): contexts for graceful shutdown
- [refctx](https://github.com/SentimensRG/ctx/tree/master/refctx): contexts linked to a reference-counter
- [mergectx](https://github.com/SentimensRG/ctx/tree/master/mergectx): utilities for merging `context.Context` instances while preserving values, errors and deadlines.

## RFC

If you find this useful please let me know:  <l.thibault@sentimens.com>

Seriously, even if you just used it in your weekend project, I'd like to hear
about it :)

## License

The MIT License

Copyright (c) 2017 Sentimens Research Group, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
