package log

import (
	"bytes"
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

func TestV(t *testing.T) {
	il, wl, el, fl := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	Info = il
	Warn = wl
	Error = el
	Fatal = fl

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
	Info = il
	Warn = wl
	Error = el
	Fatal = fl

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
	Info = il
	Warn = wl
	Error = el
	Fatal = fl

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
	Info = il
	Warn = wl
	Error = el
	Fatal = fl

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
	Info = il
	Warn = wl
	Error = el
	Fatal = fl

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
	Info = il
	Warn = wl
	Error = el
	Fatal = fl

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
	Info = il
	Warn = wl
	Error = el
	Fatal = fl

	code := -1
	Exit = func(c int) {
		code = c
	}
	ExitCode = 5

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
	if code != 5 {
		t.Errorf("Got %v, want 5 as the exit code from Fatalf", code)
	}
	if ExitCode != 5 {
		t.Errorf("ExitCode changed away from 5 after calling Fatalf: %v", ExitCode)
	}

	Exit = nil
	Fatalf("The program should not crash here")
}
