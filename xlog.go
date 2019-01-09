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
	"time"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	warnSign       = "[warn] "
	infoSign       = "[info] "
	errorSign      = "[error] "
	fatalSign      = "[fatal] "
	panicSign      = "[panic] "
)

var (
	zeroInterface  interface{}
	defaultXLogger *Logger
)

func init() {
	defaultXLogger = NewLogger(os.Stdout, "", DefaultLogFlag)
	defaultXLogger.calldepth = 3
}

type Logger struct {
	prefix    string
	flag      int
	calldepth int
	hour int
	file *os.File
	logger    *log.Logger
}

// NewLogger is similar to log.New(out io.Writer, prefix string, flag int)
func NewLogger(out io.Writer, prefix string, flag int) *Logger {
	logger := log.New(out, prefix, flag)
	l := new(Logger)
	l.prefix = prefix
	l.flag = flag
	l.logger = logger
	l.calldepth = 2
	return l
}

func NewLoggerFromFile(f *os.File, prefix string, flag int) *Logger {
	logger := log.New(f, prefix, flag)
	l := new(Logger)
	l.file = f
	l.prefix = prefix
	l.flag = flag
	l.logger = logger
	l.calldepth = 2
	l.hour = time.Now().Hour()
	return l
}

func (l *Logger) switchLogFile() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select{
		case <-ticker.C:
			//
		}
	}
}



// NewLoggerFromFileName will call os.OpenFile by os.O_CREATE|os.O_RDWR|os.O_APPEND
// and 0664 then call NewLogger()
func NewLoggerFromFileName(filename string, prefix string, flag int) *Logger {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}

	return NewLoggerFromFile(f, prefix, flag)
}

func (l *Logger) Warn(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = warnSign
	l.logger.Output(l.calldepth, fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = warnSign
	format = "%s" + format
	l.logger.Output(l.calldepth, fmt.Sprintf(format, v...))
}

func (l Logger) Warnln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = warnSign
	l.logger.Output(l.calldepth, fmt.Sprintln(v...))
}

func (l Logger) Info(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = infoSign
	l.logger.Output(l.calldepth, fmt.Sprint(v...))
}

func (l Logger) Infof(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = infoSign
	format = "%s" + format
	l.logger.Output(l.calldepth, fmt.Sprintf(format, v...))
}

func (l Logger) Infoln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = infoSign
	l.logger.Output(l.calldepth, fmt.Sprintln(v...))
}

func (l Logger) Error(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = errorSign
	l.logger.Output(l.calldepth, fmt.Sprint(v...))
}

func (l Logger) Errorf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = errorSign
	format = "%s" + format
	l.logger.Output(l.calldepth, fmt.Sprintf(format, v...))
}

func (l Logger) Errorln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = errorSign
	l.logger.Output(l.calldepth, fmt.Sprintln(v...))
}

func (l Logger) Fatal(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = fatalSign
	l.logger.Output(l.calldepth, fmt.Sprint(v...))
	os.Exit(1)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = fatalSign
	format = "%s" + format
	l.logger.Output(l.calldepth, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l Logger) Fatalln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = fatalSign
	l.logger.Output(l.calldepth, fmt.Sprintln(v...))
	os.Exit(1)
}

func (l Logger) Panic(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = panicSign

	s := fmt.Sprint(v...)
	l.logger.Output(l.calldepth, s)
	panic(s)
}

func (l Logger) Panicf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = panicSign
	format = "%s" + format
	s := fmt.Sprintf(format, v...)
	l.logger.Output(l.calldepth, s)
	panic(s)
}

func (l Logger) Panicln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = panicSign
	s := fmt.Sprintln(v...)
	l.logger.Output(l.calldepth, s)
	panic(s)
}

func Warn(v ...interface{}) {
	defaultXLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	defaultXLogger.Warnf(format, v...)
}

func Warnln(v ...interface{}) {
	defaultXLogger.Warnln(v...)
}

func Info(v ...interface{}) {
	defaultXLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	defaultXLogger.Infof(format, v...)
}

func Infoln(v ...interface{}) {
	defaultXLogger.Infoln(v...)
}

func Error(v ...interface{}) {
	defaultXLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultXLogger.Errorf(format, v...)
}

func Errorln(v ...interface{}) {
	defaultXLogger.Errorln(v...)
}

func Fatal(v ...interface{}) {
	defaultXLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultXLogger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	defaultXLogger.Fatalln(v...)
}

func Panic(v ...interface{}) {
	defaultXLogger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	defaultXLogger.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	defaultXLogger.Panicln(v...)
}
