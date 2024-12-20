package ftp

import (
	"github.com/allape/gogger"
	"github.com/fclairamb/go-log"
)

type Logger struct {
	log.Logger
	keyvals []any
	gogger  *gogger.Logger
}

func (l Logger) Debug(event string, keyvals ...interface{}) {
	l.gogger.Debug().Println(l.keyvals, event, keyvals)
}

func (l Logger) Info(event string, keyvals ...interface{}) {
	l.gogger.Info().Println(l.keyvals, event, keyvals)
}

func (l Logger) Warn(event string, keyvals ...interface{}) {
	l.gogger.Warn().Println(l.keyvals, event, keyvals)
}

func (l Logger) Error(event string, keyvals ...interface{}) {
	l.gogger.Error().Println(l.keyvals, event, keyvals)
}

func (l Logger) Panic(event string, keyvals ...interface{}) {
	l.gogger.Error().Fatalln(l.keyvals, event, keyvals)
}

func (l Logger) With(keyvals ...interface{}) log.Logger {
	return Logger{
		keyvals: keyvals,
		gogger:  l.gogger,
	}
}

func NewLogger(g *gogger.Logger) log.Logger {
	return Logger{
		gogger: g,
	}
}
