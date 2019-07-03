package sentry

import (
	"bytes"
	"io/ioutil"
	"sync"
)

type sourceReader struct {
	mu    sync.Mutex
	cache map[string][][]byte
}

func newSourceReader() sourceReader {
	return sourceReader{
		cache: make(map[string][][]byte),
	}
}

func (sr *sourceReader) readContextLines(filename string, line, context int) ([][]byte, int) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	lines, ok := sr.cache[filename]

	if !ok {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			sr.cache[filename] = nil
			return nil, 0
		}
		lines = bytes.Split(data, []byte{'\n'})
		sr.cache[filename] = lines
	}

	return sr.calculateContextLines(lines, line, context)
}

// `initial` points to a line that's the `context_line` itself in relation to returned slice
func (sr *sourceReader) calculateContextLines(lines [][]byte, line, context int) ([][]byte, int) {
	// Stacktrace lines are 1-indexed, slices are 0-indexed
	line--

	initial := context

	if lines == nil || line >= len(lines) || line < 0 {
		return nil, 0
	}

	if context < 0 {
		context = 0
		initial = 0
	}

	start := line - context

	if start < 0 {
		initial += start
		start = 0
	}

	end := line + context + 1

	if end > len(lines) {
		end = len(lines)
	}

	return lines[start:end], initial
}
