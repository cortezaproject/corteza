package zapfilter

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FilterFunc is used to check whether to filter the given entry and filters out.
type FilterFunc func(zapcore.Entry, []zapcore.Field) bool

// NewFilteringCore returns a core middleware that uses the given filter function to
// determine whether to actually call Write on the next core in the chain.
func NewFilteringCore(next zapcore.Core, filter FilterFunc) zapcore.Core {
	if filter == nil {
		filter = alwaysFalseFilter
	}
	return &filteringCore{next, filter}
}

// CheckAnyLevel determines whether at least one log level isn't filtered-out by the logger.
func CheckAnyLevel(logger *zap.Logger) bool {
	for _, level := range allLevels {
		if level >= zapcore.PanicLevel {
			continue // panic and fatal cannot be skipped
		}
		if logger.Check(level, "") != nil {
			return true
		}
	}
	return false
}

type filteringCore struct {
	next   zapcore.Core
	filter FilterFunc
}

// Check determines whether the supplied zapcore.Entry should be logged.
// If the entry should be logged, the filteringCore adds itself to the zapcore.CheckedEntry
// and returns the results.
func (core *filteringCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// FIXME: consider calling downstream core.Check too, but need to document how to
	// properly set logging level.
	if core.filter(entry, nil) {
		ce = ce.AddCore(entry, core)
	}
	return ce
}

// Write determines whether the supplied zapcore.Entry with provided []zapcore.Field should
// be logged, then calls the wrapped zapcore.Write.
func (core *filteringCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if !core.filter(entry, fields) {
		return nil
	}
	return core.next.Write(entry, fields)
}

// With adds structured context to the wrapped zapcore.Core.
func (core *filteringCore) With(fields []zapcore.Field) zapcore.Core {
	return &filteringCore{
		next:   core.next.With(fields),
		filter: core.filter,
	}
}

// Enabled asks the wrapped zapcore.Core to decide whether a given logging level is enabled
// when logging a message.
func (core *filteringCore) Enabled(level zapcore.Level) bool {
	// FIXME: Maybe it's better to always return true and only rely on the Check() func?
	//        Another way to consider it is to keep the smaller log level configured on
	//        zapfilter.
	return core.next.Enabled(level)
}

// Sync flushed buffered logs (if any).
func (core *filteringCore) Sync() error {
	return core.next.Sync()
}

// ByNamespaces takes a list of patterns to filter out logs based on their namespaces.
// Patterns are checked using path.Match.
func ByNamespaces(input string) FilterFunc {
	if input == "" {
		return alwaysFalseFilter
	}
	patterns := strings.Split(input, ",")

	// edge case optimization (always true)
	{
		hasIncludeWildcard := false
		hasExclude := false
		for _, pattern := range patterns {
			if pattern == "" {
				continue
			}
			if pattern == "*" {
				hasIncludeWildcard = true
			}
			if pattern[0] == '-' {
				hasExclude = true
			}
		}
		if hasIncludeWildcard && !hasExclude {
			return alwaysTrueFilter
		}
	}

	var mutex sync.Mutex
	matchMap := map[string]bool{}
	return func(entry zapcore.Entry, fields []zapcore.Field) bool {
		mutex.Lock()
		defer mutex.Unlock()

		if _, found := matchMap[entry.LoggerName]; !found {
			matchMap[entry.LoggerName] = false
			matchInclude := false
			matchExclude := false
			for _, pattern := range patterns {
				switch {
				case pattern[0] == '-' && !matchExclude:
					if matched, _ := path.Match(pattern[1:], entry.LoggerName); matched {
						matchExclude = true
					}
				case pattern[0] != '-' && !matchInclude:
					if matched, _ := path.Match(pattern, entry.LoggerName); matched {
						matchInclude = true
					}
				}
			}
			matchMap[entry.LoggerName] = matchInclude && !matchExclude
		}
		return matchMap[entry.LoggerName]
	}
}

// ExactLevel filters out entries with an invalid level.
func ExactLevel(level zapcore.Level) FilterFunc {
	return func(entry zapcore.Entry, fields []zapcore.Field) bool {
		return entry.Level == level
	}
}

// MinimumLevel filters out entries with a too low level.
func MinimumLevel(level zapcore.Level) FilterFunc {
	return func(entry zapcore.Entry, fields []zapcore.Field) bool {
		return entry.Level >= level
	}
}

// Any checks if any filter returns true.
func Any(filters ...FilterFunc) FilterFunc {
	return func(entry zapcore.Entry, fields []zapcore.Field) bool {
		for _, filter := range filters {
			if filter == nil {
				continue
			}
			if filter(entry, fields) {
				return true
			}
		}
		return false
	}
}

// Reverse checks is the passed filter returns false.
func Reverse(filter FilterFunc) FilterFunc {
	return func(entry zapcore.Entry, fields []zapcore.Field) bool {
		return !filter(entry, fields)
	}
}

// All checks if all filters return true.
func All(filters ...FilterFunc) FilterFunc {
	return func(entry zapcore.Entry, fields []zapcore.Field) bool {
		var atLeastOneSuccessful bool
		for _, filter := range filters {
			if filter == nil {
				continue
			}
			if !filter(entry, fields) {
				return false
			}
			atLeastOneSuccessful = true
		}
		return atLeastOneSuccessful
	}
}

// ParseRules takes a CLI-friendly set of rules to construct a filter.
//
// Syntax
//
//   pattern: RULE [RULE...]
//   RULE: one of:
//    - LEVELS:NAMESPACES
//    - NAMESPACES
//   LEVELS: LEVEL,[,LEVEL]
//   LEVEL: see `Level Patterns`
//   NAMESPACES: NAMESPACE[,NAMESPACE]
//   NAMESPACE: one of:
//    - namespace     // should be exactly this namespace
//    - *mat*ch*      // should match
//    - -NAMESPACE    // should not match
//
// Examples
//
//    *                            everything
//    *:*                          everything
//    info:*                       level info;  any namespace
//    info+:*                      levels info, warn, error, dpanic, panic, and fatal; any namespace
//    info,warn:*                  levels info, warn; any namespace
//    ns1                          any level; namespace 'ns1'
//    *:ns1                        any level; namespace 'ns1'
//    ns1*                         any level; namespaces matching 'ns1*'
//    *:ns1*                       any level; namespaces matching 'ns1*'
//    *:ns1,ns2                    any level; namespaces 'ns1' and 'ns2'
//    *:ns*,-ns3*                  any level; namespaces matching 'ns*' but not matching 'ns3*'
//    info:ns1                     level info; namespace 'ns1'
//    info,warn:ns1,ns2            levels info and warn; namespaces 'ns1' and 'ns2'
//    info:ns1 warn:n2             level info + namespace 'ns1' OR level warn and namespace 'ns2'
//    info,warn:myns* error+:*     levels info or warn and namespaces matching 'myns*' OR levels error, dpanic, panic or fatal for any namespace
func ParseRules(pattern string) (FilterFunc, error) {
	var topFilter FilterFunc

	// rules are separated by spaces, tabs or \n
	for _, rule := range strings.Fields(pattern) {
		// split rule into parts (separated by ':')
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		parts := strings.SplitN(rule, ":", 2)
		var left, right string
		switch len(parts) {
		case 1:
			// if no separator, left stays empty
			right = parts[0]
		case 2:
			if parts[0] == "" || parts[1] == "" {
				return nil, fmt.Errorf("bad syntax")
			}
			left = parts[0]
			right = parts[1]
		default:
			return nil, fmt.Errorf("bad syntax")
		}

		levelFilter, err := ByLevels(left)
		if err != nil {
			return nil, err
		}
		namespaceFilter := ByNamespaces(right)
		topFilter = Any(topFilter, All(levelFilter, namespaceFilter))
	}

	return topFilter, nil
}

// ByLevels creates a FilterFunc based on a pattern.
//
// Level Patterns
//   | Pattern | Debug | Info | Warn | Error | DPanic | Panic | Fatal |
//   | ------- | ----- | ---- | ---- | ----- | ------ | ----- | ----- |
//   | <empty> | X     | X    | X    | X     | X      | X     | X     |
//   | *       | X     | X    | X    | X     | x      | X     | X     |
//   | debug   | X     |      |      |       |        |       |       |
//   | info    |       | X    |      |       |        |       |       |
//   | warn    |       |      | X    |       |        |       |       |
//   | error   |       |      |      | X     |        |       |       |
//   | dpanic  |       |      |      |       | X      |       |       |
//   | panic   |       |      |      |       |        | X     |       |
//   | fatal   |       |      |      |       |        |       | X     |
//   | debug+  | X     | X    | x    | X     | X      | X     | X     |
//   | info+   |       | X    | X    | X     | X      | X     | X     |
//   | warn+   |       |      | X    | X     | X      | X     | X     |
//   | error+  |       |      |      | X     | X      | X     | X     |
//   | dpanic+ |       |      |      |       | X      | X     | X     |
//   | panic+  |       |      |      |       |        | X     | X     |
//   | fatal+  |       |      |      |       |        |       | X     |
func ByLevels(pattern string) (FilterFunc, error) {
	// parse pattern
	var enabled uint
	for _, part := range strings.Split(pattern, ",") {
		switch strings.ToLower(part) {
		case "", "*", "debug+":
			enabled |= debugLevel | infoLevel | warnLevel | errorLevel | dpanicLevel | panicLevel | fatalLevel
		case "debug":
			enabled |= debugLevel
		case "info":
			enabled |= infoLevel
		case "info+":
			enabled |= infoLevel | warnLevel | errorLevel | dpanicLevel | panicLevel | fatalLevel
		case "warn":
			enabled |= warnLevel
		case "warn+":
			enabled |= warnLevel | errorLevel | dpanicLevel | panicLevel | fatalLevel
		case "error":
			enabled |= errorLevel
		case "error+":
			enabled |= errorLevel | dpanicLevel | panicLevel | fatalLevel
		case "dpanic":
			enabled |= dpanicLevel
		case "dpanic+":
			enabled |= dpanicLevel | panicLevel | fatalLevel
		case "panic":
			enabled |= panicLevel
		case "panic+":
			enabled |= panicLevel | fatalLevel
		case "fatal", "fatal+":
			enabled |= fatalLevel
		default:
			return nil, fmt.Errorf("unsupported keyword: %q", pattern)
		}
	}

	// if everything is enabled
	if enabled == debugLevel&infoLevel&warnLevel&errorLevel&dpanicLevel&panicLevel&fatalLevel {
		return alwaysTrueFilter, nil
	}

	// construct custom filter
	var filter FilterFunc
	if enabled&debugLevel != 0 {
		filter = Any(ExactLevel(zapcore.DebugLevel), filter)
	}
	if enabled&infoLevel != 0 {
		filter = Any(ExactLevel(zapcore.InfoLevel), filter)
	}
	if enabled&warnLevel != 0 {
		filter = Any(ExactLevel(zapcore.WarnLevel), filter)
	}
	if enabled&errorLevel != 0 {
		filter = Any(ExactLevel(zapcore.ErrorLevel), filter)
	}
	if enabled&dpanicLevel != 0 {
		filter = Any(ExactLevel(zapcore.DPanicLevel), filter)
	}
	if enabled&panicLevel != 0 {
		filter = Any(ExactLevel(zapcore.PanicLevel), filter)
	}
	if enabled&fatalLevel != 0 {
		filter = Any(ExactLevel(zapcore.FatalLevel), filter)
	}
	return filter, nil
}

const (
	debugLevel uint = 1 << iota
	infoLevel
	warnLevel
	errorLevel
	dpanicLevel
	panicLevel
	fatalLevel
)

// MustParseRules calls ParseRules and panics if initialization failed.
func MustParseRules(pattern string) FilterFunc {
	filter, err := ParseRules(pattern)
	if err != nil {
		panic(err)
	}
	return filter
}

func alwaysFalseFilter(_ zapcore.Entry, _ []zapcore.Field) bool {
	return false
}

func alwaysTrueFilter(_ zapcore.Entry, _ []zapcore.Field) bool {
	return true
}

var allLevels = []zapcore.Level{
	zapcore.DebugLevel,
	zapcore.InfoLevel,
	zapcore.WarnLevel,
	zapcore.ErrorLevel,
	zapcore.DPanicLevel,
	zapcore.PanicLevel,
	zapcore.FatalLevel,
}
