// Copyright (c) 2015 RightScale Inc, All Rights Reserved

package log2log15

import (
	"io"
	"log"
	"strings"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
)

// NewLogger returns a stdlib log.Logger which ends up actually logging to log15. Since the stdlib
// log interface does not have a level, we need a global level to be specified here.
func NewLogger(lg log15.Logger, lvl log15.Lvl) *log.Logger {
	var w logWriter
	switch lvl {
	case log15.LvlCrit:
		w.Log = lg.Crit
	case log15.LvlError:
		w.Log = lg.Error
	case log15.LvlWarn:
		w.Log = lg.Warn
	case log15.LvlInfo:
		w.Log = lg.Info
	default:
		w.Log = lg.Debug
	}
	return log.New(w, "", 0)
}

type logWriter struct {
	Log func(msg string, ctx ...interface{})
}

func (w logWriter) Write(p []byte) (int, error) {
	// YYYY/MM/DD HH:MI:SS ...
	//var year, month,day,hour,minute,sec int
	//var msg string
	//fmt.Sscanf(string(p), "%d/%d/%d %d:%d:%d %s", &year, &month, &day, &hour, &minute, &sec, &msg)

	// strip \n
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '\n' {
			continue
		}
		if i < len(p)-1 {
			p = p[:i+1]
		}
		break
	}
	w.Log(string(p))
	return len(p), nil
}

// NewWriter returns an io.Writer suitable for use in log.New,
// and which will use the given log15.Handler
//
// This assumes that the log.Logger.Output will call the given io.Writer's
// Write method only once per log record, with the full output.
func NewWriter(handler log15.Handler) io.Writer {
	return handlerWriter{handler}
}

var _ io.Writer = handlerWriter{}

type handlerWriter struct {
	handler log15.Handler
}

func (w handlerWriter) Write(p []byte) (int, error) {
	msg := strings.TrimSpace(string(p))
	// [[YYYY/MM/DD ]HH:MI:SS[.ffffff] ][file:line: ]msg
	rec := log15.Record{Time: time.Now(), Lvl: log15.LvlDebug, Msg: msg}
	err := w.handler.Log(&rec)
	return len(p), err
}
