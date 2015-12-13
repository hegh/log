package log

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

var (
	imatcher = regexp.MustCompile("^I.*Test message\n$")
	wmatcher = regexp.MustCompile("^W.*Test message\n$")
	ematcher = regexp.MustCompile("^E.*Test message\n$")
	fmatcher = regexp.MustCompile("^F.*Test message\n$")
)

func TestLoudEnough(t *testing.T) {
	*Verbosity = 1
	if l := LoudEnough(0); !l {
		t.Errorf("Expected Verbosity=1 to be loud enough for level 0.")
	}

	if l := LoudEnough(1); !l {
		t.Errorf("Expected Verbosity=1 to be loud enough for level 1.")
	}

	if l := LoudEnough(2); l {
		t.Errorf("Expected Verbosity=1 to not be loud enough for level 2.")
	}
}

func TestCallDepth(t *testing.T) {
	// Verifies we get log_test.go (and not log.go) for the file name of log messages.
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	m := regexp.MustCompile(
		`^.*log_test\.go.*
$`)
	*Verbosity = 1

	// V
	V(1, "Test")
	if s := il.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for default V log.", s, m)
	}

	il.Truncate(0)
	Root.V(1, "Test")
	if s := il.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for Root.V log.", s, m)
	}

	// Infof
	il.Truncate(0)
	Infof("Test")
	if s := il.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for default Infof log.", s, m)
	}

	il.Truncate(0)
	Root.Infof("Test")
	if s := il.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for Root.Infof log.", s, m)
	}

	// Printf
	il.Truncate(0)
	Printf("Test")
	if s := il.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for default Printf log.", s, m)
	}

	il.Truncate(0)
	Root.Printf("Test")
	if s := il.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for Root.Printf log.", s, m)
	}

	// Warnf
	Warnf("Test")
	if s := wl.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for default Warnf log.", s, m)
	}

	wl.Truncate(0)
	Root.Warnf("Test")
	if s := wl.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for Root.Warnf log.", s, m)
	}

	// Errorf
	Errorf("Test")
	if s := el.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for default Errorf log.", s, m)
	}

	el.Truncate(0)
	Root.Errorf("Test")
	if s := el.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for Root.Errorf log.", s, m)
	}

	// Panicf
	el.Truncate(0)
	func() {
		defer func() {
			recover()
		}()
		Panicf("Test")
	}()
	if s := el.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for default Panicf log.", s, m)
	}

	el.Truncate(0)
	func() {
		defer func() {
			recover()
		}()
		Root.Panicf("Test")
	}()
	if s := el.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for Root.Panicf log.", s, m)
	}

	// Fatalf
	Root.Exit = nil
	Fatalf("Test")
	if s := fl.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for default Fatalf log.", s, m)
	}

	fl.Truncate(0)
	Root.Fatalf("Test")
	if s := fl.String(); !m.MatchString(s) {
		t.Errorf("Got %v, want something matching %v for Root.Fatalf log.", s, m)
	}
}

func TestV(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	*Verbosity = 1
	V(1, "Test %s", "message")
	V(2, "This message should not show up")
	if m := il.String(); !imatcher.MatchString(m) {
		t.Errorf("Got %v, want something matching %v from info log", m, imatcher)
	}
	if m := wl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from warn log", m)
	}
	if m := el.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from error log", m)
	}
	if m := fl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from fatal log", m)
	}
}

func TestInfo(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	Infof("Test %s", "message")
	if m := il.String(); !imatcher.MatchString(m) {
		t.Errorf("Got %v, want something matching %v from info log", m, imatcher)
	}
	if m := wl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from warn log", m)
	}
	if m := el.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from error log", m)
	}
	if m := fl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from fatal log", m)
	}
}

func TestPrint(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	Printf("Test %s", "message")
	if m := il.String(); !imatcher.MatchString(m) {
		t.Errorf("Got %v, want something matching %v from info log", m, imatcher)
	}
	if m := wl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from warn log", m)
	}
	if m := el.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from error log", m)
	}
	if m := fl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from fatal log", m)
	}
}

func TestWarn(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	Warnf("Test %s", "message")
	if m := il.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from info log", m)
	}
	if m := wl.String(); !wmatcher.MatchString(m) {
		t.Errorf("Got %v, want something matching %v from warn log", m, wmatcher)
	}
	if m := el.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from error log", m)
	}
	if m := fl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from fatal log", m)
	}
}

func TestError(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	Errorf("Test %s", "message")
	if m := il.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from info log", m)
	}
	if m := wl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from warn log", m)
	}
	if m := el.String(); !ematcher.MatchString(m) {
		t.Errorf("Got %v, want something matching %v from error log", m, ematcher)
	}
	if m := fl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from fatal log", m)
	}
}

func TestPanic(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	var err interface{}
	done := make(chan struct{})
	go func() {
		defer func() {
			err = recover()
			close(done)
		}()
		Panicf("Test %s", "message")
	}()
	<-done

	if m := il.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from info log", m)
	}
	if m := wl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from warn log", m)
	}
	if m := el.String(); !ematcher.MatchString(m) {
		t.Errorf("Got %v, want something matching %v from error log", m, ematcher)
	}
	if m := fl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from fatal log", m)
	}
}

func TestFatal(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Root.Info = il
	Root.Warn = wl
	Root.Error = el
	Root.Fatal = fl

	called := false
	Root.Exit = func() {
		called = true
	}

	Fatalf("Test %s", "message")
	if m := il.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from info log", m)
	}
	if m := wl.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from warn log", m)
	}
	if m := el.String(); len(m) > 0 {
		t.Errorf("Got %v, want empty from error log", m)
	}
	if m := fl.String(); !fmatcher.MatchString(m) {
		t.Errorf("Got %v, want something matching %v from fatal log", m, fmatcher)
	}
	if !called {
		t.Errorf("The Exit function was not called from fatal log")
	}

	Root.Exit = nil
	Fatalf("The program should not crash here")
}

type fakeTest struct {
	TestLogable
	info  *bytes.Buffer
	err   *bytes.Buffer
	fatal *bytes.Buffer
}

func (f fakeTest) Logf(format string, v ...interface{}) {
	f.info.WriteString(fmt.Sprintf(format, v...))
}

func (f fakeTest) Errorf(format string, v ...interface{}) {
	f.err.WriteString(fmt.Sprintf(format, v...))
}

func (f fakeTest) Fatalf(format string, v ...interface{}) {
	f.fatal.WriteString(fmt.Sprintf(format, v...))
}

func TestNewTest(t *testing.T) {
	// Verify NewTest() wires everything correctly for use with a test case.
	// Also verifies most of the Logger.X methods.
	ft := fakeTest{
		info:  new(bytes.Buffer),
		err:   new(bytes.Buffer),
		fatal: new(bytes.Buffer),
	}
	lg := NewTest(ft, "TestNewTest", true)

	lg.Infof("Info log")
	lg.Printf("Print log")
	lg.Warnf("Warn log")
	lg.Errorf("Error log")
	lg.Fatalf("Fatal log")

	info := regexp.MustCompile(
		`^I\d{2}:\d{2}:\d{2}\.\d{6} log_test.go:\d+: Info log
I.*Print log
W.*Warn log
$`)
	if s := ft.info.String(); !info.MatchString(s) {
		t.Errorf("Got %v, want something matching %v from info log", s, info)
	}

	err := regexp.MustCompile(
		`^E\d{2}:\d{2}:\d{2}\.\d{6} log_test.go:\d+: Error log
$`)
	if s := ft.err.String(); !err.MatchString(s) {
		t.Errorf("Got %v, want something matching %v from error log", s, err)
	}

	fatal := regexp.MustCompile(
		`^F\d{2}:\d{2}:\d{2}\.\d{6} log_test.go:\d+: Fatal log
$`)
	if s := ft.fatal.String(); !fatal.MatchString(s) {
		t.Errorf("Got %v, want something matching %v from fatal log", s, fatal)
	}

	// Reset and verify failOnError=false directs Errorf() to the info log.
	ft.info.Truncate(0)
	lg = NewTest(ft, "TestNewTest", false)
	lg.Errorf("Error log")
	if s := ft.err.String(); !err.MatchString(s) {
		t.Errorf("Got %v, want something matching %v from error log", s, err)
	}
}
