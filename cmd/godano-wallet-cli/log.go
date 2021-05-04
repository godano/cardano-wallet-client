package main

import (
	"bytes"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// By setting GODANO_WALLET_CLIENT_VERBOSE, allow debug-level output before the main Cobra command is executed
func (c *walletCLI) configureEarlyLogLevel() {
	if os.Getenv("GODANO_WALLET_CLIENT_VERBOSE") != "" {
		c.log.SetLevel(logrus.DebugLevel)
	}
}

func (c *walletCLI) configureLogLevel() {
	level := logrus.InfoLevel
	if c.logVerbose {
		level = logrus.TraceLevel
	} else if c.logVerbose {
		level = logrus.DebugLevel
	} else if c.logVeryQuiet {
		level = logrus.ErrorLevel
	} else if c.logQuiet {
		level = logrus.WarnLevel
	}
	c.log.SetLevel(level)
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
	// Avoid empty log lines due to trailing \n characters
	text = bytes.TrimSpace(text)
	text = append(text, '\n')
	return text, nil
}
