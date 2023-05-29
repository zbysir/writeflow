package log

import "testing"

func TestColor(t *testing.T) {
	//SetDev(true)
	buf := BuffSink{}
	l := New(Options{
		IsDev:         true,
		To:            &buf,
		DisableCaller: true,
		CallerSkip:    0,
		Name:          "",
	})

	l.Infof("%v", "info")
	l.Debugf("%v", "debug")
	l.Warnf("%v", "warn")
	l.Errorf("%v", "error")

	t.Logf("buf: %s", buf.buf.String())
}

func TestFormat(t *testing.T) {
	l := New(Options{
		IsDev:         true,
		DisableCaller: true,
		CallerSkip:    0,
		Name:          "",
	})

	l.Infof("123")
}
