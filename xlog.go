// Package log implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. Logger accessible through helper functions
// Warn[f|ln], Info[f,ln], Error[f,ln], Fatal[f|ln], and Panic[f|ln], which are
// easier to use than creating a Logger manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// Every log message is output on a separate line: if the message being
// printed does not end in a newline, the logger will add one.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.

package xlog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

const (
	// DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	levelWarn      = "warn"
	levelInfo      = "info"
	levelError     = "error"
	levelFatal     = "fatal"
	levelPanic     = "panic"
)

var (
	zeroInterface     interface{}
	defaultXLogger    *Logger
	defaultBufferSize = 1024
)

func init() {
	defaultXLogger = NewLogger(os.Stdout, "", DefaultLogFlag)
	defaultXLogger.calldepth = 3
}

type logContent struct {
	t      time.Time
	level  string
	file   string
	line   int
	format string
	v      []interface{}
}

type Logger struct {
	prefix    string
	flag      int
	calldepth int
	out       io.Writer
	buffer    chan logContent
}

// time | level | file | msg
func (l Logger) format(lc logContent) []byte {
	buf := bytes.Buffer{}
	year, month, day := lc.t.Date()
	hour, min, sec := lc.t.Clock()

	buf.WriteString(strconv.Itoa(year))
	buf.WriteByte('/')
	if month < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.Itoa(int(month)))
	buf.WriteByte('/')
	if day < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.Itoa(day))
	buf.WriteByte(' ')
	if hour < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.Itoa(hour))
	buf.WriteByte(':')
	if min < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.Itoa(min))
	buf.WriteByte(':')
	if sec < 10 {
		buf.WriteByte('0')
	}
	buf.WriteString(strconv.Itoa(sec))
	buf.WriteByte(' ')
	buf.WriteByte('[')
	buf.WriteString(lc.level)
	buf.WriteByte(']')
	buf.WriteByte(' ')
	buf.WriteString(lc.file)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(lc.line))
	buf.WriteByte(' ')
	if lc.format == "" {
		buf.WriteString(fmt.Sprint(lc.v...))
	} else {
		buf.WriteString(fmt.Sprintf(lc.format, lc.v...))
	}

	return buf.Bytes()
}

// NewLogger is similar to log.New(out io.Writer, prefix string, flag int)
func NewLogger(out io.Writer, prefix string, flag int) *Logger {
	l := new(Logger)
	l.prefix = prefix
	l.flag = flag
	l.out = out
	l.calldepth = 3
	l.buffer = make(chan logContent, defaultBufferSize)

	go l.write()
	return l
}

func NewLoggerFromFile(logFile string, prefix string, flag int) *Logger {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return NewLogger(f, prefix, flag)
}

func (l *Logger) write() {
	for {
		select {
		case lc := <-l.buffer:
			logBytes := l.format(lc)
			l.out.Write(logBytes)
		}
	}
}

func (l *Logger) AddCallDepth(n int) {
	l.calldepth += n
}

func (l *Logger) output(level string, format string, v ...interface{}) {
	t := time.Now()
	_, file, line, ok := runtime.Caller(l.calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	l.buffer <- logContent{
		t:      t,
		level:  level,
		file:   file,
		line:   line,
		format: format,
		v:      v,
	}
}

func (l *Logger) Warn(v ...interface{}) {
	l.Warnf("", v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(levelWarn, format, v...)
}

func (l Logger) Info(v ...interface{}) {
	l.Infof("", v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	l.output(levelInfo, format, v...)
}

func (l Logger) Error(v ...interface{}) {
	l.Errorf("", v...)
}

func (l Logger) Errorf(format string, v ...interface{}) {
	l.output(levelError, format, v...)
}

func (l Logger) Fatal(v ...interface{}) {
	l.Fatalf("", v...)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.output(levelFatal, format, v...)
	os.Exit(1)
}

func (l Logger) Panic(v ...interface{}) {
	l.Panicf("", v...)
}

func (l Logger) Panicf(format string, v ...interface{}) {
	l.output(levelPanic, format, v...)

	var s string
	if format == "" {
		s = fmt.Sprint(v...)
	} else {
		s = fmt.Sprintf(format, v...)
	}
	panic(s)
}

func Warn(v ...interface{}) {
	defaultXLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	defaultXLogger.Warnf(format, v...)
}

func Info(v ...interface{}) {
	defaultXLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	defaultXLogger.Infof(format, v...)
}

func Error(v ...interface{}) {
	defaultXLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultXLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	defaultXLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultXLogger.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	defaultXLogger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	defaultXLogger.Panicf(format, v...)
}
