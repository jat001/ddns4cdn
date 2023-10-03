package core

import (
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

type (
	LogEntry  = logrus.Entry
	LogFields = logrus.Fields
)

type log struct {
	Logger *logrus.Logger
}

func (ctx *log) SetFormatter() {
	ctx.Logger.SetFormatter(&nested.Formatter{
		TimestampFormat: time.DateTime,
		ShowFullLevel:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"module", "submodule"},
	})
}

func (ctx *log) SetLevel(level string) {
	l, e := logrus.ParseLevel(level)
	if e != nil {
		l = logrus.DebugLevel
	}
	ctx.Logger.SetLevel(l)
}

var Log = log{
	Logger: logrus.New(),
}
