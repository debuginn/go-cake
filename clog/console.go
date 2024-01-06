package clog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	ConsoleLogInitFailed = errors.New("init XConsoleLog failed, not found level")
)

type ConsoleLog struct {
	level    int
	skip     int
	hostname string
	service  string
}

type Brush func(string) string

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;37"), // white
	NewBrush("1;36"), // debug cyan
	NewBrush("1;33"), // info  green
	NewBrush("1;31"), // warn  yellow
	NewBrush("1;32"), // fatal red
}

func init() {
	_ = RegisterLogger(Console, NewConsoleLog())
}

func NewConsoleLog() LogInterface {
	return &ConsoleLog{
		skip: LogDefSkipNum,
	}
}

func (c *ConsoleLog) Init(config *LoggerConfig) error {

	if config == nil {
		return ConsoleLogInitFailed
	}

	c.level = LevelFromStr(config.Level)
	c.service = config.Service
	if config.Skip > 0 {
		c.skip = config.Skip
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	c.hostname = hostname

	return nil
}

func (c *ConsoleLog) Reopen() error {
	return nil
}

func (c *ConsoleLog) SetLevel(level string) {
	c.level = LevelFromStr(level)
}

func (c *ConsoleLog) SetSkip(skip int) {
	c.skip = skip
}

func (c *ConsoleLog) Fatal(ctx context.Context, format string, a ...interface{}) {

	if c.level > FatalLevel {
		return
	}

	logText := Format(format, a...)
	fun, filename, lineno := GetRuntimeInfo(c.skip)

	color := colors[FatalLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	c.write(ctx, FatalLevel, &logText)
}

func (c *ConsoleLog) Warn(ctx context.Context, format string, a ...interface{}) {

	if c.level > WarnLevel {
		return
	}

	logText := Format(format, a...)
	fun, filename, lineno := GetRuntimeInfo(c.skip)

	color := colors[WarnLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	c.write(ctx, WarnLevel, &logText)
}

func (c *ConsoleLog) Info(ctx context.Context, format string, a ...interface{}) {

	if c.level > InfoLevel {
		return
	}

	logText := Format(format, a...)
	fun, filename, lineno := GetRuntimeInfo(c.skip)

	color := colors[InfoLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	c.write(ctx, InfoLevel, &logText)
}

func (c *ConsoleLog) Debug(ctx context.Context, format string, a ...interface{}) {

	if c.level > DebugLevel {
		return
	}

	logText := Format(format, a...)
	fun, filename, lineno := GetRuntimeInfo(c.skip)

	color := colors[DebugLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	c.write(ctx, DebugLevel, &logText)
}

func (c *ConsoleLog) Close() {

}

func (c *ConsoleLog) write(ctx context.Context, level int, msg *string) {

	color := colors[level]
	levelText := color(levelTextArray[level])
	t := time.Now().Format("2006-01-02 15:04:05")

	traceId := GetTraceId(ctx)

	logText := FormatLog(msg, t, c.service, c.hostname, levelText, traceId)
	file := os.Stdout
	if level >= WarnLevel {
		file = os.Stderr
	}

	_, _ = file.Write([]byte(logText))
}
