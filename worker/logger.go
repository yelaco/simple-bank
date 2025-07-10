package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Print(level zerolog.Level, args ...any) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...any) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

// Debug logs a message at Debug level.
func (logger *Logger) Debug(args ...any) {
	logger.Print(zerolog.DebugLevel, args...)
}

// Info logs a message at Info level.
func (logger *Logger) Info(args ...any) {
	logger.Print(zerolog.InfoLevel, args...)
}

// Warn logs a message at Warning level.
func (logger *Logger) Warn(args ...any) {
	logger.Print(zerolog.WarnLevel, args...)
}

// Error logs a message at Error level.
func (logger *Logger) Error(args ...any) {
	logger.Print(zerolog.ErrorLevel, args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (logger *Logger) Fatal(args ...any) {
	logger.Print(zerolog.FatalLevel, args...)
}
