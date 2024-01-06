package clog

import "context"

// LogInterface .
type LogInterface interface {

	// Init  log instance
	Init(config *LoggerConfig) error
	// Reopen log reopen
	Reopen() error
	// SetLevel set logger print level
	SetLevel(level string)
	// SetSkip set skip runtime info
	SetSkip(skip int)

	Fatal(ctx context.Context, format string, a ...interface{})
	Warn(ctx context.Context, format string, a ...interface{})
	Info(ctx context.Context, format string, a ...interface{})
	Debug(ctx context.Context, format string, a ...interface{})

	Close()
}
