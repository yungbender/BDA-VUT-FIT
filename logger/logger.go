package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Prefix string
}

func (l *Logger) Info(fields logrus.Fields, msg string) {
	logrus.WithFields(fields).Info(l.Prefix + ": " + msg)
}

func (l *Logger) Warning(fields logrus.Fields, msg string) {
	logrus.WithFields(fields).Warning(l.Prefix + ": " + msg)
}

func (l *Logger) Error(fields logrus.Fields, msg string) {
	logrus.WithFields(fields).Error(l.Prefix + ": " + msg)
}
