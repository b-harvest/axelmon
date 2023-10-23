package log

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var eventQueue chan func()

func init() {
	// Logger setup
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC1123,
	}
	log.Logger = log.Output(output)

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.ErrorStackFieldName = "trace"

	eventQueue = make(chan func())

	// For thread safe
	go func() {
		for event := range eventQueue {
			event()
		}
	}()
}

func enqueue(event func()) {
	eventQueue <- event
}

func Info(msg string) {
	event := func() {
		log.Info().Msg(msg)
	}
	enqueue(event)
}

func Error(err error) {
	stack := string(debug.Stack())
	event := func() {
		log.Error().Err(err).Msg("\n" + stack)
	}
	enqueue(event)
}

func Debug(msg any) {
	message := fmt.Sprint(msg)
	event := func() {
		log.Debug().Msg(message)
	}
	enqueue(event)
}
