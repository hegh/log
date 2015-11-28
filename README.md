A basic log package that sits on top of Go's existing log package.

Provides logging levels (Info, Warning, Error, Fatal), verbosity,  and a couple
of side effects (log-and-panic, log-and-terminate). Also makes redirecting log
output easier.

Use like you would the existing logging package:
    import "github.com/hegh/log"

    func f() {
      x := 5
      log.Infof("This is an info-level message. x = %v", x)
    }

Redirect logging output by setting `log.Info`, `log.Warn`, `log.Error`, and
`log.Fatal` to alternative io.Writer instances. Use io.MultiWriter for more
interesting setups (like having separate Fatal, Error, Warning, and Info log
files, with each lower-level file receiving all of the more severe messages
in addition to its own).
