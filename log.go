// Package log provides a helpful wrapper around the standard log package.
//
// Anticipated basic usage:
// log.Infof("This is an info level message")
// log.Warnf("This is a warn level message")
// log.Errorf("This is an error level message")
// log.V(5, "This is info level, but will only show up if --verbosity >= 5")
// log.Panicf("This message is error level, and also becomes a panic()")
// log.Fatalf("This message is fatal level, and os.Exit(1) follows immediately")
package log

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// Kept separate because these may diverge.
	// Call depth when using the default calls to the Root Logger.
	rootDepth = 3

	// Call depth when directly calling log functions on a Logger.
	logDepth = 3
)

var (
	Verbosity = flag.Int("verbosity", 0, "Logging verbosity level. Higher means more logs.")
	Root      *Logger
)

// The rewriter type allows us to change the destination of written data without
// rebuilding the actual log.Logger objects used.
type rewriter struct {
	w *io.Writer
}

func (w *rewriter) Write(p []byte) (int, error) {
	return (*w.w).Write(p)
}

func init() {
	Root = New("")
}

// Logable is the interface required for writing data to the next lower level.
type Logable interface {
	// Output a log message. See log.Logger.Output for details.
	Output(calldepth int, s string) error
}

// Logger provides an individually configurable logging instance.
type Logger struct {
	name string

	// Verbosity indicates how "loud" this logger is.
	// It defaults to the Verbosity flag.
	Verbosity *int

	i, w, e, f Logable

	// Info is where all INFO-level messages get written.
	Info io.Writer

	// Warn is where all WARN-level messages get written.
	Warn io.Writer

	// Error is where all ERROR-level messages (including Panic) get written.
	Error io.Writer

	// Fatal is where all FATAL-level messages get written.
	Fatal io.Writer

	// Exit is the function to call after logging a Fatal message.
	// If nil, is not called.
	Exit func()
}

// New returns a new Logger with the given name.
func New(name string) *Logger {
	l := &Logger{
		name:      name,
		Verbosity: Verbosity,
		Info:      os.Stderr,
		Warn:      os.Stderr,
		Error:     os.Stderr,
		Fatal:     os.Stderr,
		Exit:      func() { os.Exit(1) },
	}
	flags := log.Ldate | log.Ltime | log.Lshortfile
	l.i = log.New(&rewriter{&l.Info}, "I", flags)
	l.w = log.New(&rewriter{&l.Warn}, "W", flags)
	l.e = log.New(&rewriter{&l.Error}, "E", flags)
	l.f = log.New(&rewriter{&l.Fatal}, "F", flags)
	return l
}

func (l *Logger) Name() string {
	return l.name
}

// SetVerbosity is a convenience method to set the logging verbosity to a constant.
func (l *Logger) SetVerbosity(v int) {
	l.Verbosity = &v
}

// Formats the message and writes it to the given logger.
// Returns the formatted message.
// If there is an error writing to the given logger, writes a description
// including the given message to the base logger.
func write(l Logable, depth int, name, format string, v ...interface{}) string {
	msg := fmt.Sprintf(format, v...)
	if err := l.Output(depth, msg); err != nil {
		log.Printf("Failed to write to %s logger: %v.\n  Message: %s", name, err, msg)
	}
	return msg
}

// LoudEnough returns whether the verbosity is high enough to include messages of the given level.
func (l *Logger) LoudEnough(level int) bool {
	return level <= *l.Verbosity
}

// LoudEnough returns whether the verbosity on the root logger is high enough to include messages of the given level.
func LoudEnough(level int) bool {
	return Root.LoudEnough(level)
}

// V writes log messages at INFO level, but only if the configured verbosity is equal or greater than the provided level.
func (l *Logger) V(level int, format string, v ...interface{}) {
	if l.LoudEnough(level) {
		write(l.i, logDepth, l.name+" info", format, v...)
	}
}

// V writes log messages at INFO level to the root logger, but only if the configured verbosity is equal or greater than the provided level.
func V(level int, format string, v ...interface{}) {
	if Root.LoudEnough(level) {
		write(Root.i, rootDepth, Root.name+" info", format, v...)
	}
}

// Infof writes log messages at INFO level.
func (l *Logger) Infof(format string, v ...interface{}) {
	write(l.i, logDepth, l.name+" info", format, v...)
}

// Infof writes log messages at INFO level to the root logger.
func Infof(format string, v ...interface{}) {
	write(Root.i, rootDepth, Root.name+" info", format, v...)
}

// Printf is synonymous with Infof.
// It exists for compatibility with the basic log package.
func (l *Logger) Printf(format string, v ...interface{}) {
	write(l.i, logDepth, l.name+" info", format, v...)
}

// Printf is synonymous with Infof.
// It exists for compatibility with the basic log package.
func Printf(format string, v ...interface{}) {
	write(Root.i, rootDepth, Root.name+" info", format, v...)
}

// Warnf writes log messages at WARN level.
func (l *Logger) Warnf(format string, v ...interface{}) {
	write(l.w, logDepth, l.name+" warn", format, v...)
}

// Warnf writes log messages at WARN level to the root logger.
func Warnf(format string, v ...interface{}) {
	write(Root.w, rootDepth, Root.name+" warn", format, v...)
}

// Errorf writes log messages at ERROR level.
func (l *Logger) Errorf(format string, v ...interface{}) {
	write(l.e, logDepth, l.name+" error", format, v...)
}

// Errorf writes log messages at ERROR level to the root logger.
func Errorf(format string, v ...interface{}) {
	write(Root.e, rootDepth, Root.name+" error", format, v...)
}

// Panicf writes log messages at ERROR level, and then panics.
// The panic parameter is an error with the formatted message.
func (l *Logger) Panicf(format string, v ...interface{}) {
	panic(errors.New(write(l.e, logDepth, l.name+" error", format, v...)))
}

// Panicf writes log messages at ERROR level to the root logger, and then panics.
// The panic parameter is an error with the formatted message.
func Panicf(format string, v ...interface{}) {
	panic(errors.New(write(Root.e, rootDepth, Root.name+" error", format, v...)))
}

// Fatalf writes log messages at FATAL level, and then calls Exit.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	write(l.f, logDepth, l.name+" fatal", format, v...)
	if l.Exit != nil {
		l.Exit()
	}
}

// Fatalf writes log messages at FATAL level to the root logger, and then calls Exit.
func Fatalf(format string, v ...interface{}) {
	write(Root.f, rootDepth, Root.name+" fatal", format, v...)
	if Root.Exit != nil {
		Root.Exit()
	}
}
