package wallet

import (
	"bytes"
	"time"

	"github.com/sirupsen/logrus"
)

// Package-internal logger. Can be exported when required.
var Log = logrus.New()

func init() {
	formatter := newLogFormatter()
	logrus.StandardLogger().SetFormatter(formatter)
	Log.SetFormatter(formatter)
}

func newLogFormatter() *myFormatter {
	return &myFormatter{
		f: logrus.TextFormatter{
			DisableColors:   false,
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.StampMilli,
		},
	}
}

type myFormatter struct {
	f logrus.TextFormatter
}

func (f *myFormatter) Format(e *logrus.Entry) ([]byte, error) {
	text, err := f.f.Format(e)
	if err != nil {
		return text, err
	}
	// Remove all whitespace and replace with a single trailing newline character
	// Many libraries explicitly add a \n character to log lines, which leads to empty lines.
	text = bytes.TrimSpace(text)
	text = append(text, '\n')
	return text, nil
}
