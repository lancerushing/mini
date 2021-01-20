package logutil

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configure zerolog.
func Configure() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if os.Getenv("GAE_ENV") != "" {
		log.Logger = log.With().Logger().
			Level(zerolog.InfoLevel).
			Hook(AddStackDriverSeverity{}).
			Hook(StackDriverSourceLocation).
			Output(os.Stdout)
	} else {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
