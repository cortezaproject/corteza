[![Tests on Linux, MacOS and Windows](https://github.com/bep/godartsass/workflows/Test/badge.svg)](https://github.com/bep/godartsass/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/godartsass)](https://goreportcard.com/report/github.com/bep/godartsass)
[![codecov](https://codecov.io/gh/bep/godartsass/branch/main/graph/badge.svg?token=OWZ9RCAYWO)](https://codecov.io/gh/bep/godartsass)
[![GoDoc](https://godoc.org/github.com/bep/godartsass?status.svg)](https://godoc.org/github.com/bep/godartsass)

This is a Go API backed by the native [Dart Sass](https://github.com/sass/dart-sass/releases) executable running with `sass --embedded`.

>**Note:** The `v2.x.x` of this project targets the `v2` of the Dart Sass Embedded protocol with the `sass` exexutable in releases that can be downloaeded [here](https://github.com/sass/dart-sass/releases). For `v1` you need to import `github.com/bep/godartsass` and not `github.com/bep/godartsass/v2`.

The primary motivation for this project is to provide `SCSS` support to [Hugo](https://gohugo.io/). I welcome PRs with bug fixes. I will also consider adding functionality, but please raise an issue discussing it first.

For LibSass bindings in Go, see [GoLibSass](https://github.com/bep/golibsass).

```

