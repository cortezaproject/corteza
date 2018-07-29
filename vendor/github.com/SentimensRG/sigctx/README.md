# DEPRECATION WARNING

`sigctx` is now part of [SentimensRG/ctx](https://github.com/SentimensRG/ctx).

This repository will no longer receive updates.

# sigctx
Go contexts for graceful shutdown

## installation

```bash
go get -u github.com/SentimensRG/sigctx
```

## usage

```go
ctx := sigctx.New()  // returns a regular context.Context

// With this simple pattern, your goroutines are guaranteed to terminate correctly
ctx, cancel := context.WithCancel(ctx)
go someBlockingFunction(ctx)
defer cancel()

<-ctx.Done()  // will unblock on SIGINT and SIGTERM
```

## RFC

If you find this useful, please let me know:  <l.thibault@sentimens.com>

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
