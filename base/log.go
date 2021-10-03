package base

import (
	log "github.com/sirupsen/logrus"
)

func init() {

}

var Log logs

type logs struct {
}

func (l *logs) Info(args ...interface{}) {
	log.Info(args)
}

func (l *logs) Warn(args ...interface{}) {
	log.Warn(args)
}

func (l *logs) Fatal(args ...interface{}) {
	log.Fatal(args)
}

func (l *logs) Error(args ...interface{}) {
	log.Error(args)
}
