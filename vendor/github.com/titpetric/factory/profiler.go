package factory

import (
	"fmt"
	"sync"
	"time"
)

// DatabaseProfiler is an interface to provide log/timing data for issued SQL queries
type DatabaseProfiler interface {
	Post(*DatabaseProfilerContext)
	Flush()
}

// DatabaseProfilerContext contains the issued SQL query and timestamp
type DatabaseProfilerContext struct {
	Query string
	Args  string
	Time  time.Time
}

// new is the private constructor for a DatabaseProfilerContext
func (DatabaseProfilerContext) new(query string, args ...interface{}) *DatabaseProfilerContext {
	a := fmt.Sprintf("%#v", args)
	return &DatabaseProfilerContext{
		Query: query,
		Args:  a[15 : len(a)-1], // omits `[]interface {...}`
		Time:  time.Now(),
	}
}

// duration returns the time in seconds since the context was created
func (p *DatabaseProfilerContext) duration() float64 {
	return time.Since(p.Time).Seconds()
}

// String returns the info about the query (fmt.Stringer interface)
func (p *DatabaseProfilerContext) String() string {
	return fmt.Sprintf("[%.4fs] %s (%s)", p.duration(), p.Query, p.Args)
}

// DatabaseProfilerStdout logs query statistics to stdout (fmt.Printf)
type DatabaseProfilerStdout struct {
}

// Post prints the query statistics to stdout
func (*DatabaseProfilerStdout) Post(p *DatabaseProfilerContext) {
	fmt.Println(p)
}

// Flush stdout (no-op for this profiler)
func (*DatabaseProfilerStdout) Flush() {
}

// DatabaseProfilerMemory logs query statistics to internal buffer
type DatabaseProfilerMemory struct {
	sync.Mutex

	Log []string
}

// Post prints the query statistics
func (this *DatabaseProfilerMemory) Post(p *DatabaseProfilerContext) {
	this.Lock()
	this.Log = append(this.Log, p.String())
	this.Unlock()
}

// Flush log to stdout with fmt.Println
func (this *DatabaseProfilerMemory) Flush() {
	this.Lock()
	count := len(this.Log)
	for _, line := range this.Log[:count] {
		fmt.Println(line)
	}
	this.Log = this.Log[count:]
	this.Unlock()
}
