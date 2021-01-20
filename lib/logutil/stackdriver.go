package logutil

import (
	"encoding/json"
	"log"
	"runtime"
	"strings"

	"cloud.google.com/go/logging"
	"github.com/rs/zerolog"
	logging2 "google.golang.org/api/logging/v2"
)

// AddStackDriverSeverity will add a "severity" with stackdriver appropriate values.
type AddStackDriverSeverity struct{}

// Run the hook.
func (h AddStackDriverSeverity) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		severity := logging.Default

		switch level {
		case zerolog.DebugLevel:
			severity = logging.Debug
		case zerolog.InfoLevel:
			severity = logging.Info
		case zerolog.WarnLevel:
			severity = logging.Warning
		case zerolog.ErrorLevel:
			severity = logging.Error
		case zerolog.FatalLevel:
			severity = logging.Critical
		case zerolog.PanicLevel:
			severity = logging.Critical
		default:
			severity = logging.Warning
		}

		e.Str("severity", severity.String())
	}
}

var skipCount = zerolog.CallerSkipFrameCount + 1

// StackDriverSourceLocation instance of type sdSourceLocation struct {.
var StackDriverSourceLocation = newSDSourceLocation(skipCount)

func newSDSourceLocation(skipFrameCount int) sdSourceLocation {
	return sdSourceLocation{callerSkipFrameCount: skipFrameCount}
}

type sdSourceLocation struct {
	callerSkipFrameCount int
}

func (ch sdSourceLocation) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	pc, file, line, ok := runtime.Caller(ch.callerSkipFrameCount)

	if !ok {
		log.Print("no caller")

		return
	}

	if strings.Contains(file, "github.com/jackc/pgx/v4/log/zerologadapter/adapter.go") {
		pc, file, line, _ = runtime.Caller(ch.callerSkipFrameCount + 1)
	}

	details := runtime.FuncForPC(pc)
	s := logging2.LogEntrySourceLocation{
		File:     file,
		Function: details.Name(),
		Line:     int64(line),
	}

	b, err := json.Marshal(s)
	if err != nil {
		log.Print(err)

		return
	}

	e.RawJSON("sourceLocation", b)
}
