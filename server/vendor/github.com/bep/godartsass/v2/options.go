package godartsass

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/godartsass/v2/internal/embeddedsass"
)

// Options configures a Transpiler.
type Options struct {
	// The path to the Dart Sass wrapper binary, an absolute filename
	// if not in $PATH.
	// If this is not set, we will try 'dart-sass'
	// (or 'dart-sass.bat' on Windows) in the OS $PATH.
	// There may be several ways to install this, one would be to
	// download it from here: https://github.com/sass/dart-sass/releases
	DartSassEmbeddedFilename string

	// Timeout is the duration allowed for dart sass to transpile.
	// This was added for the beta6 version of Dart Sass Protocol,
	// as running this code against the beta5 binary would hang
	// on Execute.
	Timeout time.Duration

	// LogEventHandler will, if set, receive log events from Dart Sass,
	// e.g. @debug and @warn log statements.
	LogEventHandler func(LogEvent)
}

// LogEvent is a type of log event from Dart Sass.
type LogEventType int

const (
	// Usually triggered by the @warn directive.
	LogEventTypeWarning LogEventType = iota

	// Events trigered for usage of deprecated Sass features.
	LogEventTypeDeprecated

	// Triggered by the @debug directive.
	LogEventTypeDebug
)

type LogEvent struct {
	// Type is the type of log event.
	Type LogEventType

	// Message on the form url:line:col message.
	Message string
}

func (opts *Options) init() error {
	if opts.DartSassEmbeddedFilename == "" {
		opts.DartSassEmbeddedFilename = defaultDartSassBinaryFilename
	}

	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	return nil
}

// ImportResolver allows custom import resolution.
//
// CanonicalizeURL should create a canonical version of the given URL if it's
// able to resolve it, else return an empty string.
//
// A canonicalized URL should include a scheme, e.g. 'file:///foo/bar.scss',
// if applicable, see:
//
//	https://en.wikipedia.org/wiki/File_URI_scheme
//
// Importers   must ensure that the same canonical URL
// always refers to the same stylesheet.
//
// Load loads the canonicalized URL's content.
type ImportResolver interface {
	CanonicalizeURL(url string) (string, error)
	Load(canonicalizedURL string) (Import, error)
}

type Import struct {
	// The content of the imported file.
	Content string

	// The syntax of the imported file.
	SourceSyntax SourceSyntax
}

// Args holds the arguments to Execute.
type Args struct {
	// The input source.
	Source string

	// The URL of the Source.
	// Leave empty if it's unknown.
	// Must include a scheme, e.g. 'file:///myproject/main.scss'
	// See https://en.wikipedia.org/wiki/File_URI_scheme
	//
	// Note: There is an open issue for this value when combined with custom
	// importers, see https://github.com/sass/dart-sass/issues/24
	URL string

	// Defaults is SCSS.
	SourceSyntax SourceSyntax

	// Default is EXPANDED.
	OutputStyle OutputStyle

	// If enabled, a sourcemap will be generated and returned in Result.
	EnableSourceMap bool

	// If enabled, sources will be embedded in the generated source map.
	SourceMapIncludeSources bool

	// Custom resolver to use to resolve imports.
	// If set, this will be the first in the resolver chain.
	ImportResolver ImportResolver

	// Additional file paths to uses to resolve imports.
	IncludePaths []string

	sassOutputStyle  embeddedsass.OutputStyle
	sassSourceSyntax embeddedsass.Syntax

	// Ordered list starting with options.ImportResolver, then IncludePaths.
	sassImporters []*embeddedsass.InboundMessage_CompileRequest_Importer
}

func (args *Args) init(seq uint32, opts Options) error {
	if args.OutputStyle == "" {
		args.OutputStyle = OutputStyleExpanded
	}
	if args.SourceSyntax == "" {
		args.SourceSyntax = SourceSyntaxSCSS
	}

	v, ok := embeddedsass.OutputStyle_value[string(args.OutputStyle)]
	if !ok {
		return fmt.Errorf("invalid OutputStyle %q", args.OutputStyle)
	}
	args.sassOutputStyle = embeddedsass.OutputStyle(v)

	v, ok = embeddedsass.Syntax_value[string(args.SourceSyntax)]
	if !ok {
		return fmt.Errorf("invalid SourceSyntax %q", args.SourceSyntax)
	}

	args.sassSourceSyntax = embeddedsass.Syntax(v)

	if args.ImportResolver != nil {
		args.sassImporters = []*embeddedsass.InboundMessage_CompileRequest_Importer{
			{
				Importer: &embeddedsass.InboundMessage_CompileRequest_Importer_ImporterId{
					ImporterId: seq,
				},
			},
		}
	}

	if args.IncludePaths != nil {
		for _, p := range args.IncludePaths {
			args.sassImporters = append(args.sassImporters, &embeddedsass.InboundMessage_CompileRequest_Importer{Importer: &embeddedsass.InboundMessage_CompileRequest_Importer_Path{
				Path: filepath.Clean(p),
			}})
		}
	}

	return nil
}

type (
	// OutputStyle defines the style of the generated CSS.
	OutputStyle string

	// SourceSyntax defines the syntax of the source passed in Execute.
	SourceSyntax string
)

const (
	// Expanded (default) output.
	// Note that LibSASS and Ruby SASS have more output styles, and their
	// default is NESTED.
	OutputStyleExpanded OutputStyle = "EXPANDED"

	// Compressed/minified output.
	OutputStyleCompressed OutputStyle = "COMPRESSED"
)

const (
	// SCSS style source syntax (default).
	SourceSyntaxSCSS SourceSyntax = "SCSS"

	// Indented or SASS style source syntax.
	SourceSyntaxSASS SourceSyntax = "INDENTED"

	// Regular CSS source syntax.
	SourceSyntaxCSS SourceSyntax = "CSS"
)

// ParseOutputStyle will convert s into OutputStyle.
// Case insensitive, returns OutputStyleNested for unknown value.
func ParseOutputStyle(s string) OutputStyle {
	switch OutputStyle(strings.ToUpper(s)) {
	case OutputStyleCompressed:
		return OutputStyleCompressed
	case OutputStyleExpanded:
		return OutputStyleExpanded
	default:
		return OutputStyleExpanded
	}
}

// ParseSourceSyntax will convert s into SourceSyntax.
// Case insensitive, returns SourceSyntaxSCSS for unknown value.
func ParseSourceSyntax(s string) SourceSyntax {
	switch SourceSyntax(strings.ToUpper(s)) {
	case SourceSyntaxSCSS:
		return SourceSyntaxSCSS
	case SourceSyntaxSASS, "SASS":
		return SourceSyntaxSASS
	case SourceSyntaxCSS:
		return SourceSyntaxCSS
	default:
		return SourceSyntaxSCSS
	}
}
