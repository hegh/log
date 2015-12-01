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

var (
	Verbosity = flag.Int("verbosity", 0, "Logging verbosity level. Higher means more logs.")

	// Info is where all INFO-level messages get written.
	Info io.Writer = os.Stderr

	// Warn is where all WARN-level messages get written.
	Warn io.Writer = os.Stderr

	// Error is where all ERROR-level messages (including Panic) get written.
	Error io.Writer = os.Stderr

	// Fatal is where all FATAL-level messages get written.
	Fatal io.Writer = os.Stderr
)

// The rewriter type allows us to change the destination of written data without
// rebuilding the actual log.Logger objects used.
type rewriter struct {
	w *io.Writer
}

func (w *rewriter) Write(p []byte) (int, error) {
	return (*w.w).Write(p)
}

var (
	// The loggers used internally.
	i, w, e, f *log.Logger
)

func init() {
	flags := log.Ldate | log.Ltime | log.Lshortfile
	i = log.New(&rewriter{&Info}, "I", flags)
	w = log.New(&rewriter{&Warn}, "W", flags)
	e = log.New(&rewriter{&Error}, "E", flags)
	f = log.New(&rewriter{&Fatal}, "F", flags)
}

// Formats the message and writes it to the given logger.
// Returns the formatted message.
// If there is an error writing to the given logger, writes a description
// including the given message to the base logger.
func write(l *log.Logger, name, format string, v ...interface{}) string {
	msg := fmt.Sprintf(format, v...)
	if err := l.Output(3, msg); err != nil {
		log.Printf("Failed to write to %s logger: %v.\n  Message: %s", name, err, msg)
	}
	return msg
}

// LoudEnough returns whether the verbosity is high enough to include messages of the given level.
func LoudEnough(level int) bool {
	return level <= *Verbosity
}

// V writes log messages at INFO level, but only if the configured verbosity is equal or greater than the provided level.
func V(level int, format string, v ...interface{}) {
	if LoudEnough(level) {
		write(i, "info", format, v...)
	}
}

// Infof writes log messages at INFO level.
func Infof(format string, v ...interface{}) {
	write(i, "info", format, v...)
}

// Printf is synonymous with Infof.
// It exists for compatibility with the basic log package.
func Printf(format string, v ...interface{}) {
	write(i, "info", format, v...)
}

// Warnf writes log messages at WARN level.
func Warnf(format string, v ...interface{}) {
	write(w, "warn", format, v...)
}

// Errorf writes log messages at ERROR level.
func Errorf(format string, v ...interface{}) {
	write(e, "error", format, v...)
}

// Panicf writes log messages at ERROR level, and then panics.
// The panic parameter is an error with the formatted message.
func Panicf(format string, v ...interface{}) {
	panic(errors.New(write(e, "error", format, v...)))
}

// Fatalf writes log messages at FATAL level, and then calls os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	write(f, "fatal", format, v...)
	os.Exit(1)
}
