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
	"sync"
	"time"
)

const (
	// DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	levelDebug     = "debug"
	levelWarn      = "warn"
	levelInfo      = "info"
	levelError     = "error"
	levelFatal     = "fatal"
	levelPanic     = "panic"
)

//LogLevel log level
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var (
	zeroInterface     interface{}
	defaultXLogger    *Logger
	defaultBufferSize = 1024
)

func init() {
	defaultXLogger = NewLogger(os.Stdout, Options{})
	defaultXLogger.calldepth = 3
}

type logContent struct {
	t      time.Time
	level  LogLevel
	file   string
	line   int
	format string
	v      []interface{}
}

type Logger struct {
	level     LogLevel
	prefix    string
	flag      int
	calldepth int
	out       io.Writer
	f         *os.File
	buffer    chan logContent

	fileName   string
	dayChange  chan bool
	curDay     int
	bufferPool *sync.Pool
}

// time | level | file | msg
func (l Logger) format(lc logContent) []byte {
	buf := l.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer l.bufferPool.Put(buf)

	year, month, day := lc.t.Date()
	hour, min, sec := lc.t.Clock()

	if day != l.curDay {
		l.curDay = day
		l.dayChange <- true
	}

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

	switch lc.level {
	case LevelDebug:
		buf.WriteString(levelDebug)
	case LevelInfo:
		buf.WriteString(levelInfo)
	case LevelWarn:
		buf.WriteString(levelWarn)
	case LevelError:
		buf.WriteString(levelError)
	case LevelPanic:
		buf.WriteString(levelPanic)
	case LevelFatal:
		buf.WriteString(levelFatal)
	}

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
	buf.WriteByte('\n')

	return buf.Bytes()
}

type Options struct {
	Prefix string
	Level  LogLevel
}

// NewLogger is similar to log.New(out io.Writer, prefix string, flag int)
func NewLogger(out io.Writer, opts Options) *Logger {
	l := new(Logger)
	l.prefix = opts.Prefix
	// l.flag = flag
	l.level = opts.Level
	l.out = out
	l.calldepth = 3
	l.buffer = make(chan logContent, defaultBufferSize)
	l.bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	l.curDay = time.Now().Day()
	l.dayChange = make(chan bool)
	go l.write()
	go l.changeFileByDay()
	return l
}

func NewLoggerFromFile(logFile string, opts Options) *Logger {
	nowLogFile := logFile + "." + formatTime(time.Now())
	f, err := createFile(nowLogFile)
	if err != nil {
		log.Fatal(err)
	}

	l := NewLogger(f, opts)
	l.fileName = logFile
	l.f = f

	if err := os.Symlink(nowLogFile, logFile); err != nil {
		log.Fatal(err)
	}

	return l
}

func (l *Logger) AddCalldepth(num int) {
	l.calldepth += num
}

func (l *Logger) changeFileByDay() {
	for {
		select {
		case ok := <-l.dayChange:
			if ok && l.f != nil {

				// 新建一个文件
				nowLogFile := l.fileName + "." + formatTime(time.Now())
				f, err := createFile(nowLogFile)
				if err != nil {
					l.Error(err)
					continue
				}
				oldF := l.f
				l.out = f
				l.f = f

				if oldF != nil {
					oldF.Close()
				}
				// 建立连接
				if err := os.Symlink(nowLogFile, l.fileName); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%4d%2d%2d", t.Year(), t.Month(), t.Day())
}

func createFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
}

func (l *Logger) write() {
	for {
		select {
		case lc := <-l.buffer:
			logBytes := l.format(lc)
			l.out.Write(logBytes)
			if lc.level == LevelFatal {
				os.Exit(1)
			} else if lc.level == LevelPanic {
				var s string
				if lc.format == "" {
					s = fmt.Sprint(lc.v...)
				} else {
					s = fmt.Sprintf(lc.format, lc.v...)
				}
				panic(s)
			}
		}
	}
}

func (l *Logger) output(level LogLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}

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

func (l *Logger) SetLogLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) Debug(v ...interface{}) {
	l.Debugf("", v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(LevelDebug, format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Warnf("", v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(LevelWarn, format, v...)
}

func (l Logger) Info(v ...interface{}) {
	l.Infof("", v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	l.output(LevelInfo, format, v...)
}

func (l Logger) Error(v ...interface{}) {
	l.Errorf("", v...)
}

func (l Logger) Errorf(format string, v ...interface{}) {
	l.output(LevelError, format, v...)
}

func (l Logger) Fatal(v ...interface{}) {
	l.Fatalf("", v...)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.output(LevelFatal, format, v...)
}

func (l Logger) Panic(v ...interface{}) {
	l.Panicf("", v...)
}

func (l Logger) Panicf(format string, v ...interface{}) {
	l.output(LevelPanic, format, v...)
}

func Debug(v ...interface{}) {
	defaultXLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	defaultXLogger.Debugf(format, v...)
}

func Warn(v ...interface{}) {
	defaultXLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	defaultXLogger.Warnf(format, v...)
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
