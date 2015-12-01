A basic log package that sits on top of Go's existing log package.

Provides logging levels (Info, Warning, Error, Fatal), verbosity,  and a couple
of side effects (log-and-panic, log-and-terminate). Also makes redirecting log
output easier.

Use generally like you would the existing logging package:
    import "github.com/hegh/log"

    func f() {
      x := 5
      log.Infof("This is an info-level message. x = %v", x)
      log.Printf("This is also info-level, provided for API compatibility")
    }

Redirect logging output by setting `log.Info`, `log.Warn`, `log.Error`, and
`log.Fatal` to alternative `io.Writer` instances.

Here is an example of complex log redirection. This directs everything into
`progname.info.log`, warning and above into `progname.warn.log`, error and above
into `progname.error.log`, and fatal messages additionally get printed to
`stderr`:

    // Package preinit should be imported by any package that needs to log
    // anything during its initialization.
    // Just `import _ "preinit"` and use log.XXX as normal.
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
      log.Info = i
      log.Warn = io.MultiWriter(log.Info, w)
      log.Error = io.MultiWriter(log.Warn, e)
      log.Fatal = io.MultiWriter(log.Error, os.Stderr)
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
