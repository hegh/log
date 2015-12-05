A basic log package that sits on top of Go's existing log package, but also
allows log redirection to unit test logs.

# Main features

 * Logging levels (Info, Warning, Error, Fatal).
 * Control over verbosity for debug logs.
 * Log-and-panic.
 * Log-and-call-a-function (by default os.Exit(1)).
 * Runtime redirection of log output.

# Basic usage

Use generally like you would the existing logging package:

    import "github.com/hegh/log"

    func f() {
      x := 5
      log.Infof("This is an info-level message. x = %v", x)
      log.Printf("This is also info-level, provided for API compatibility")
    }

Redirect logging output by setting `log.Root.Info`, `log.Root.Warn`,
`log.Root.Error`, and `log.Root.Fatal` to alternative `io.Writer` instances.

# Advanced usage

Here is an example of complex log redirection. This directs everything into
`progname.info.log`, warning and above into `progname.warn.log`, error and above
into `progname.error.log`, and fatal messages additionally get printed to
`stderr`:

    // Package preinit should be imported by any package that needs to log
    // anything during initialization.
    // Just `import _ "preinit"` in addition to `log` and use log.XXX as normal.
    package preinit

    import (
      "flag"
      "io"
      "os"
    )

    var (
      i, w, e *os.File
    )

    func init() {
      var err error
      if i, err = os.Create(flag.Arg(0) + ".info.log"); err != nil {
        panic(err)
      }
      if w, err = os.Create(flag.Arg(0) + ".warn.log"); err != nil {
        panic(err)
      }
      if e, err  = os.Create(flag.Arg(0) + "error.log"); err != nil {
        panic(err)
      }
      log.Root.Info = i
      log.Root.Warn = io.MultiWriter(log.Root.Info, w)
      log.Root.Error = io.MultiWriter(log.Root.Warn, e)
      log.Root.Fatal = io.MultiWriter(log.Root.Error, os.Stderr)
    }

    // Close should be called prior to program termination.
    // Probably best to `defer preinit.Close()` in `func main()`.
    func Close() {
      // Not a bad idea to check the errors from these and return them...
      i.Close()
      w.Close()
      e.Close()
    }

And then:

    package main

    import (
      "log"
      "preinit"
    )

    func main() {
      defer preinit.Close()

      // Do stuff. Log it. Enjoy.
    }

# Unit test usage

Assuming single-threaded unit tests (no use of `t.Parallel()`), just replace the
Root logger at the beginning of each test function:

    import (
      "testing"

      "github.com/hegh/log"
    )

    func TestSomethingCool(t *testing.T) {
      log.Root = log.NewTest(t, "TestSomethingCool", false)
      log.Infof("Only printed if the test fails (or is run with -v). Has an 'I' indicator.")
      log.Warnf("Same rules as Infof, but with a 'W'.")
      log.Errorf("Also like Infof, but with an 'E'.")
      log.Fatalf("The test fails and stops here, with an 'F'.")
    }

Now, info, warning, and error messages go to `t.Logf`. If you want error
messages to go to `t.Errorf` and cause the test to fail, pass `true` instead of
`false` to `log.NewTest()`. Fatal messages will always go to `t.Fatalf` (causing
the test to fail and immediately abort).

