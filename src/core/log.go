package core

import (
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

type log struct {
	SetFormatter func()
	SetLevel     func(string)
}

func setFormatter() {
	Logger.SetFormatter(&nested.Formatter{
		TimestampFormat: time.DateTime,
		ShowFullLevel:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"module", "submodule"},
	})
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
