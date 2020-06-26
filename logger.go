// Copyright (c) 2012-present The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package db

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type LogLevel int8

const (
	LogLevelTrace LogLevel = -1

	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelPanic
)

var logLevels = map[LogLevel]string{
	LogLevelTrace: "TRACE",
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARNING",
	LogLevelError: "ERROR",
	LogLevelFatal: "FATAL",
	LogLevelPanic: "PANIC",
}

const (
	defaultLogLevel LogLevel = LogLevelWarn
)

var defaultLogger Logger = log.New(os.Stdout, "", log.LstdFlags)

// Logger
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Print(v ...interface{})
	Printf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

// LoggingCollector represents a logging collector.
type LoggingCollector interface {
	Enabled(LogLevel) bool

	SetLogger(Logger)
	SetLevel(LogLevel)

	Trace(v ...interface{})
	Tracef(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

type loggingCollector struct {
	level  LogLevel
	logger Logger
}

func (c *loggingCollector) Enabled(level LogLevel) bool {
	return level >= c.level
}

func (c *loggingCollector) SetLevel(level LogLevel) {
	c.level = level
}

func (c *loggingCollector) Level() LogLevel {
	return c.level
}

func (c *loggingCollector) Logger() Logger {
	if c.logger == nil {
		return defaultLogger
	}
	return c.logger
}

func (c *loggingCollector) SetLogger(logger Logger) {
	c.logger = logger
}

func (c *loggingCollector) logf(level LogLevel, f string, v ...interface{}) {
	format := logLevels[level] + "\n" + f
	if _, file, line, ok := runtime.Caller(2); ok {
		format = fmt.Sprintf("log_level=%s file=%s:%d\n%s", logLevels[level], file, line, f)
	}
	format = "upper/db: " + format

	if level >= LogLevelPanic {
		c.Logger().Panicf(format, v...)
	}
	if level >= LogLevelFatal {
		c.Logger().Fatalf(format, v...)
	}
	if c.Enabled(level) {
		c.Logger().Printf(format, v...)
	}
}

func (c *loggingCollector) log(level LogLevel, v ...interface{}) {
	format := logLevels[level] + "\n"
	if _, file, line, ok := runtime.Caller(2); ok {
		format = fmt.Sprintf("log_level=%s file=%s:%d\n", logLevels[level], file, line)
	}
	format = "upper/db: " + format
	v = append([]interface{}{format}, v...)

	if level >= LogLevelPanic {
		c.Logger().Panic(v...)
	}
	if level >= LogLevelFatal {
		c.Logger().Fatal(v...)
	}
	if c.Enabled(level) {
		c.Logger().Print(v...)
	}
}

func (c *loggingCollector) Debugf(format string, v ...interface{}) {
	c.logf(LogLevelDebug, format, v...)
}
func (c *loggingCollector) Debug(v ...interface{}) {
	c.log(LogLevelDebug, v...)
}

func (c *loggingCollector) Tracef(format string, v ...interface{}) {
	c.logf(LogLevelTrace, format, v...)
}
func (c *loggingCollector) Trace(v ...interface{}) {
	c.log(LogLevelDebug, v...)
}

func (c *loggingCollector) Infof(format string, v ...interface{}) {
	c.logf(LogLevelInfo, format, v...)
}
func (c *loggingCollector) Info(v ...interface{}) {
	c.log(LogLevelInfo, v...)
}

func (c *loggingCollector) Warnf(format string, v ...interface{}) {
	c.logf(LogLevelWarn, format, v...)
}
func (c *loggingCollector) Warn(v ...interface{}) {
	c.log(LogLevelWarn, v...)
}

func (c *loggingCollector) Errorf(format string, v ...interface{}) {
	c.logf(LogLevelError, format, v...)
}
func (c *loggingCollector) Error(v ...interface{}) {
	c.log(LogLevelError, v...)
}

func (c *loggingCollector) Fatalf(format string, v ...interface{}) {
	c.logf(LogLevelFatal, format, v...)
}
func (c *loggingCollector) Fatal(v ...interface{}) {
	c.log(LogLevelFatal, v...)
}

func (c *loggingCollector) Panicf(format string, v ...interface{}) {
	c.logf(LogLevelPanic, format, v...)
}
func (c *loggingCollector) Panic(v ...interface{}) {
	c.log(LogLevelPanic, v...)
}

var defaultLoggingCollector LoggingCollector = &loggingCollector{
	level:  defaultLogLevel,
	logger: defaultLogger,
}

func Log() LoggingCollector {
	return defaultLoggingCollector
}

func init() {
	if logLevel := os.Getenv("UPPER_DB_LOG"); logLevel != "" {
		for k := range logLevels {
			if logLevels[k] == logLevel {
				Log().SetLevel(k)
				return
			}
		}
	}
}
