package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

var errorlog *slog.Logger

type Logger struct {
	*slog.Logger
}

func Error(msg string, args ...any) {
	errorlog.Error(msg, args...)
}
func (l *Logger) Error(msg string, args ...any) {

	stack := make([]uintptr, 10)
	length := runtime.Callers(2, stack)
	stack = stack[:length]

	errstr := ""

	for i, pc := range stack {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)

		errstr += fmt.Sprintf("%d. %s:%d  ", i+1, file, line)
	}
	args = append(args, "Stack")
	args = append(args, errstr)
	errorlog.Error(msg, args...)

}

var filemap map[string]*os.File = make(map[string]*os.File)

var logmap map[string]*Logger = make(map[string]*Logger)

// func Init(basedir string, level slog.Level) {
// 	logfile := basedir + "logs/slog.log"

// 	dir := filepath.Dir(logfile)
// 	err := os.MkdirAll(dir, os.ModePerm)
// 	if err != nil {
// 		panic(err)
// 	}
// 	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	filemap["default"] = file

// 	slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(file, os.Stdout), &slog.HandlerOptions{AddSource: false, Level: level})))

// 	dir = fmt.Sprintf("%slogs/", basedir)
// 	//logfile := "/"
// 	//dir := filepath.Dir(basedir + logfile)
// 	err = os.MkdirAll(dir, os.ModePerm)
// 	if err != nil {
// 		panic(err)
// 	}
// 	file, err = os.OpenFile(dir+"error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	filemap["errorlog"] = file
// 	errorlog = slog.New(slog.NewTextHandler(io.MultiWriter(file, os.Stdout), &slog.HandlerOptions{AddSource: true, Level: slog.LevelError}))

// }

var once sync.Once

func Close() error {
	for _, v := range filemap {
		v.Close()
	}
	return nil
}

func GetLog(name string, ty string, basedir string) *Logger {

	if _, ok := logmap[name]; !ok {
		dir := fmt.Sprintf("%slogs/%s/%s/", basedir, ty, name)
		//logfile := "/"
		//dir := filepath.Dir(basedir + logfile)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
		// 获取当前日期
		currentDate := time.Now().Format("20060102")

		// 拼接文件名
		fileName := dir + "slog_" + currentDate + ".log"
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		filemap[name] = file
		logmap[name] = &Logger{slog.New(slog.NewTextHandler(io.MultiWriter(file, os.Stdout),
			&slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}))}
	}
	once.Do(func() {

		file, err := os.OpenFile(basedir+"logs/errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		filemap["errors"] = file
		errorlog = slog.New(slog.NewTextHandler(io.MultiWriter(file, os.Stdout), &slog.HandlerOptions{AddSource: false, Level: slog.LevelError}))
	})
	return logmap[name]
}

// type LogLevel int

// const (
// 	DEBUG LogLevel = iota
// 	INFO
// 	WARNING
// 	ERROR
// )

// type Logger struct {
// 	mu            sync.Mutex
// 	logLevel      LogLevel
// 	logDir        string
// 	logFile       *os.File
// 	infoLogger    *log.Logger
// 	warningLogger *log.Logger
// 	errorLogger   *log.Logger
// }

// func NewLogger(logDir string, logLevel LogLevel) (*Logger, error) {
// 	if _, err := os.Stat(logDir); os.IsNotExist(err) {
// 		if err := os.MkdirAll(logDir, 0755); err != nil {
// 			return nil, err
// 		}
// 	}
// 	logFile, err := os.OpenFile(filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02"))), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logger := &Logger{
// 		logLevel:      logLevel,
// 		logDir:        logDir,
// 		logFile:       logFile,
// 		infoLogger:    log.New(io.MultiWriter(os.Stdout, logFile), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
// 		warningLogger: log.New(io.MultiWriter(os.Stdout, logFile), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
// 		errorLogger:   log.New(io.MultiWriter(os.Stderr, logFile), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
// 	}
// 	return logger, nil
// }
// func (l *Logger) SetLogLevel(level LogLevel) {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()
// 	l.logLevel = level
// }
// func (l *Logger) Info(v ...interface{}) {
// 	if l.logLevel > INFO {
// 		return
// 	}
// 	l.infoLogger.Output(2, fmt.Sprintln(v...))
// }
// func (l *Logger) Warning(v ...interface{}) {
// 	if l.logLevel > WARNING {
// 		return
// 	}
// 	l.warningLogger.Output(2, fmt.Sprintln(v...))
// }
// func (l *Logger) Error(v ...interface{}) {
// 	if l.logLevel > ERROR {
// 		return
// 	}
// 	l.errorLogger.Output(2, fmt.Sprintln(v...))
// }
// func (l *Logger) Close() error {
// 	return l.logFile.Close()
// }
