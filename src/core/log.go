package core

import (
	"github.com/sirupsen/logrus"
)

type log struct {
	SetFormatter func()
	SetLevel     func(string)
}

func setFormatter() {
	Logger.SetFormatter(&logrus.TextFormatter{})
}

func setLevel(level string) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		l = logrus.DebugLevel
	}
	Logger.SetLevel(l)
}

type (
	LogEntry  = *logrus.Entry
	LogFields = logrus.Fields
)

var (
	Log = log{
		SetFormatter: setFormatter,
		SetLevel:     setLevel,
	}
	Logger = logrus.New()
)
