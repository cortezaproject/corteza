package logger

import (
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	rrBufEntry struct {
		num int
		rec []byte
	}

	// simple round-robin struct that holds our buffer with log entries
	rr struct {
		mux sync.RWMutex
		num int
		buf []*rrBufEntry
	}

	debugBufferingLogger struct {
		zapcore.LevelEnabler
		out *rr
		enc zapcore.Encoder
	}
)

const (
	debugLogCap = 10240
)

var (
	// initialize debug logger round robin db
	debugLogRR = &rr{
		buf: make([]*rrBufEntry, 0),
	}
)

func WriteLogBuffer(w io.Writer, after, limit int) (_ int, err error) {
	var (
		// was at least one entry outputted?
		has bool
	)

	debugLogRR.mux.RLock()
	defer debugLogRR.mux.RUnlock()
	if _, err = w.Write([]byte{'['}); err != nil {
		return
	}

	for _, e := range debugLogRR.buf {
		if after >= e.num {
			continue
		}

		// count back to zero
		limit--

		if has {
			if _, err = w.Write([]byte{','}); err != nil {
				return
			}
		}

		if _, err = w.Write(e.rec); err != nil {
			return
		}

		if limit == 0 {
			break
		}

		has = true
	}

	if _, err = w.Write([]byte{']'}); err != nil {
		return
	}

	return
}

func (r *rr) append(ent []byte) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.num++

	// modify serialized json
	var bufEnt = &rrBufEntry{r.num, append(ent[:len(ent)-2], []byte(fmt.Sprintf(`,"index":%d}`, r.num))...)}

	if len(r.buf) >= debugLogCap {
		r.buf = append(r.buf[1:], bufEnt)
	} else {
		r.buf = append(r.buf, bufEnt)
	}
}

func DebugBufferedLogger(out *rr) *debugBufferingLogger {
	var encConf = zap.NewProductionEncoderConfig()
	encConf.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return &debugBufferingLogger{
		LevelEnabler: zapcore.DebugLevel,
		out:          out,
		enc:          zapcore.NewJSONEncoder(encConf),
	}
}

func (c *debugBufferingLogger) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (c *debugBufferingLogger) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *debugBufferingLogger) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	encbuf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	c.out.append(encbuf.Bytes())
	return nil
}

func (c *debugBufferingLogger) Sync() error { return nil }
func (c *debugBufferingLogger) clone() *debugBufferingLogger {
	return &debugBufferingLogger{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc.Clone(),
		out:          c.out,
	}
}
