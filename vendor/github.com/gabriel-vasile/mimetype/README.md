<h1 align="center">
  mimetype
</h1>

<h4 align="center">
  A package for detecting MIME types and extensions based on magic numbers
</h4>
<h6 align="center">
  No C bindings, zero dependencies and thread safe
</h6>

<p align="center">
  <a href="https://travis-ci.org/gabriel-vasile/mimetype">
    <img alt="Build Status" src="https://travis-ci.org/gabriel-vasile/mimetype.svg?branch=master">
  </a>
  <a href="https://godoc.org/github.com/gabriel-vasile/mimetype">
    <img alt="Documentation" src="https://godoc.org/github.com/gabriel-vasile/mimetype?status.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/gabriel-vasile/mimetype">
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/gabriel-vasile/mimetype">
  </a>
  <a href="https://coveralls.io/github/gabriel-vasile/mimetype?branch=master">
    <img alt="Go report card" src="https://coveralls.io/repos/github/gabriel-vasile/mimetype/badge.svg?branch=master">
  </a>
  <a href="LICENSE">
    <img alt="License" src="https://img.shields.io/badge/License-MIT-green.svg">
  </a>
</p>

## Features
- fast and precise MIME type and file extension detection
- long list of [supported MIME types](supported_mimes.md)
- common file formats are prioritized
- small and simple API
- handles MIME type aliases
- thread safe
- low memory usage, besides the file header

## Install
```bash
go get github.com/gabriel-vasile/mimetype
```

## Usage
There are quick [examples](EXAMPLES.md) and
[GoDoc](https://godoc.org/github.com/gabriel-vasile/mimetype) for full reference.


## Structure
**mimetype** uses an hierarchical structure to keep the MIME type detection logic.
This reduces the number of calls needed for detecting the file type. The reason
behind this choice is that there are file formats used as containers for other
file formats. For example, Microsoft Office files are just zip archives,
containing specific metadata files. Once a file a file has been identified as a
zip, there is no need to check if it is a text file, but it is worth checking if
it is an Microsoft Office file.

To prevent loading entire files into memory, when detecting from a
[reader](https://godoc.org/github.com/gabriel-vasile/mimetype#DetectReader)
or from a [file](https://godoc.org/github.com/gabriel-vasile/mimetype#DetectFile)
**mimetype** limits itself to reading only the header of the input.
<div align="center">
  <img alt="structure" src="https://github.com/gabriel-vasile/mimetype/blob/33abbe6cb78fe1a8486c92f95008a9e0fcef10a1/mimetype.gif?raw=true" width="88%">
</div>

## Performance
Thanks to the hierarchical structure, searching for common formats first,
and limiting itself to file headers, **mimetype** matches the performance of
stdlib `http.DetectContentType` while outperforming the alternative package.

[Benchmarks](https://github.com/gabriel-vasile/mimetype/blob/d8628c314b5e59259afc7b0f4f84e6b31931b316/mimetype_test.go#L267)
were run on an Intel Xeon Gold 6136 24 core CPU @ 3.00GHz. Lower is better.
```bash
                            mimetype  http.DetectContentType      filetype
BenchmarkMatchTar-24       250 ns/op         400 ns/op           3778 ns/op
BenchmarkMatchZip-24       524 ns/op         351 ns/op           4884 ns/op
BenchmarkMatchJpeg-24      103 ns/op         228 ns/op            839 ns/op
BenchmarkMatchGif-24       139 ns/op         202 ns/op            751 ns/op
BenchmarkMatchPng-24       165 ns/op         221 ns/op           1176 ns/op
```

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md).
