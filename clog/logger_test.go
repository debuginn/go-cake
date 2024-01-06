package clog

import (
	"context"
	"fmt"
	"testing"
)

func TestConsoleLogger(t *testing.T) {

	config := &LoggerConfig{
		Level:   LevelDebug,
		Service: "console_service",
	}

	err := initLogger(Console, config, "")
	if err != nil {
		_ = fmt.Errorf("TestConsoleLogger initLogger failed, err:%v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, TraceId, "asdsdsdsd")

	Debug(ctx, "console debug test, %d", 1)
	Info(ctx, "console info test, %d", 1)
	Warn(ctx, "console warn test, %d", 1)
	Fatal(ctx, "console fatal test, %d", 1)
}
