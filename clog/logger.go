package clog

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// LogInstance .
type LogInstance struct {
	logger  LogInterface
	enable  bool
	initial bool
	source  string
}

type LoggerConfig struct {
	Level   string
	Source  string
	Service string
	Skip    int
}

var (
	gLoggerMgr map[string]*LogInstance = make(map[string]*LogInstance)
	lock       sync.RWMutex
)

func init() {
	config := &LoggerConfig{
		Level: LevelDebug,
	}
	_ = initLogger(Console, config, SourceAuto)
}

func initLogger(name string, config *LoggerConfig, source string) error {

	// Check whether logs are implemented
	instance, ok := gLoggerMgr[name]
	if !ok {
		err := errors.New(fmt.Sprintf("not found logger[%s]", name))
		return err
	}

	// init logger
	err := instance.logger.Init(config)
	if err != nil {
		return err
	}

	if len(source) > 0 {
		instance.source = source
	} else {
		instance.source = ""
	}

	instance.enable = true
	instance.initial = true

	return nil
}

func InitLogger(name string, config *LoggerConfig) error {

	lock.Lock()
	defer lock.Unlock()

	err := initLogger(name, config, "")
	if err != nil {
		return err
	}

	// close auto inject logger
	for _, v := range gLoggerMgr {
		if v.logger == nil || !v.enable {
			continue
		}

		if v.source == SourceAuto {
			v.enable = false
		}
	}

	return nil
}

func RegisterLogger(name string, logger LogInterface) error {

	lock.Lock()
	defer lock.Unlock()

	_, ok := gLoggerMgr[name]
	if ok {
		err := errors.New(fmt.Sprintf("duplicate logger: %s", name))
		return err
	}

	gLoggerMgr[name] = &LogInstance{
		logger:  logger,
		enable:  false,
		initial: false,
	}

	return nil
}

func EnableLogger(name string, enable bool) error {

	lock.Lock()
	defer lock.Unlock()

	instance, ok := gLoggerMgr[name]
	if !ok {
		err := errors.New(fmt.Sprintf("not found logger: %s", name))
		return err
	}

	if !instance.initial {
		instance.enable = false
		return nil
	}

	instance.enable = true
	return nil
}

func GetLogger(name string) (logger LogInterface, err error) {

	lock.Lock()
	defer lock.Unlock()

	instance, ok := gLoggerMgr[name]
	if !ok {
		err := errors.New(fmt.Sprintf("not found logger: %s", name))
		return nil, err
	}

	return instance.logger, nil
}

func UnregisterLogger(name string) error {

	lock.Lock()
	defer lock.Unlock()

	instance, ok := gLoggerMgr[name]
	if !ok {
		err := errors.New(fmt.Sprintf("not found logger: %s", name))
		return err
	}

	if instance != nil {
		instance.logger.Close()
	}

	delete(gLoggerMgr, name)
	return nil
}

func Reopen() error {

	lock.Lock()
	defer lock.Unlock()

	var errMsg string

	for k, v := range gLoggerMgr {

		// logger is nil or logger not enable
		if v.logger == nil || !v.enable {
			continue
		}

		rerr := v.logger.Reopen()
		if rerr != nil {
			errMsg += fmt.Sprintf("logger: %s reload failed, err:%+v \n", k, rerr)
			continue
		}
	}

	err := errors.New(errMsg)
	return err
}

func SetLevelAll(level string) {

	lock.Lock()
	defer lock.Unlock()

	for _, v := range gLoggerMgr {

		// logger is nil or logger not enable
		if v.logger == nil || !v.enable {
			continue
		}

		v.logger.SetLevel(level)
	}
}

func SetLevel(name, level string) error {

	lock.Lock()
	defer lock.Unlock()

	instance, ok := gLoggerMgr[name]
	if !ok {
		err := errors.New(fmt.Sprintf("not found logger: %s", name))
		return err
	}

	instance.logger.SetLevel(level)
	return nil
}

func Fatal(ctx context.Context, format string, a ...interface{}) {

	lock.RLock()
	defer lock.RUnlock()

	for _, v := range gLoggerMgr {
		// logger is nil or logger not enable
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Fatal(ctx, format, a...)
	}

	return
}

func Warn(ctx context.Context, format string, a ...interface{}) {

	lock.RLock()
	defer lock.RUnlock()

	for _, v := range gLoggerMgr {
		// logger is nil or logger not enable
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Warn(ctx, format, a...)
	}

	return
}

func Info(ctx context.Context, format string, a ...interface{}) {

	lock.RLock()
	defer lock.RUnlock()

	for _, v := range gLoggerMgr {
		// logger is nil or logger not enable
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Info(ctx, format, a...)
	}

	return
}

func Debug(ctx context.Context, format string, a ...interface{}) {

	lock.RLock()
	defer lock.RUnlock()

	for _, v := range gLoggerMgr {
		// logger is nil or logger not enable
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Debug(ctx, format, a...)
	}

	return
}

func Close() {

	lock.RLock()
	defer lock.RUnlock()

	for _, v := range gLoggerMgr {
		// logger is nil or logger not enable
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Close()
	}

	return
}
