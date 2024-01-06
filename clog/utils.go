package clog

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	FatalLevel
	NoneLevel
)

const (
	TraceId       = "traceId"
	LogDefSkipNum = 6

	Console = "console"

	SourceAuto = "auto"

	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelFatal = "fatal"
	LevelNone  = "none"
)

var (
	levelTextArray = []string{
		DebugLevel: "DEBUG",
		InfoLevel:  "INFO",
		WarnLevel:  "WARN",
		FatalLevel: "FATAL",
	}
)

func LevelFromStr(level string) int {

	resultLevel := DebugLevel
	levelLower := strings.ToLower(level)
	switch levelLower {
	case LevelDebug:
		resultLevel = DebugLevel
	case LevelInfo:
		resultLevel = InfoLevel
	case LevelWarn:
		resultLevel = WarnLevel
	case LevelFatal:
		resultLevel = FatalLevel
	case LevelNone:
		resultLevel = NoneLevel
	default:
		resultLevel = InfoLevel
	}

	return resultLevel
}

func Format(format string, a ...interface{}) (result string) {

	if len(a) == 0 {
		result = format
		return
	}

	return fmt.Sprintf(format, a...)
}

func GetRuntimeInfo(skip int) (function, filename string, lineno int) {

	function = "???"
	pc, filename, lineno, ok := runtime.Caller(skip)
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}

	return function, filename, lineno
}

func FormatLog(body *string, fields ...string) string {

	var buffer bytes.Buffer
	for _, v := range fields {
		buffer.WriteString("[")
		buffer.WriteString(v)
		buffer.WriteString("] ")
	}

	buffer.WriteString(*body)
	buffer.WriteString("\n")

	return buffer.String()
}

func GetTraceId(ctx context.Context) string {

	if ctx == nil {
		return ""
	}

	if v := ctx.Value(TraceId); v != nil {
		if vi, ok := v.(int); ok {
			return strconv.Itoa(vi)
		} else if vi, ok := v.(int64); ok {
			return strconv.FormatInt(vi, 10)
		} else if vs, ok := v.(string); ok {
			return vs
		} else {
			return fmt.Sprintf("%v", v)
		}
	}

	return ""
}
