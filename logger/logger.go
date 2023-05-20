package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	WriteTrace(string)
	WriteError(string)
}

type ZeroLogger struct {
	ConsoleLogger
	FileLogger
}

type FileLogger struct {
	Logger zerolog.Logger
	File   *os.File
}

type ConsoleLogger struct {
	Logger zerolog.Logger
}

func NewLogger() Logger {
	return &ZeroLogger{
		ConsoleLogger: ConsoleLogger{Logger: InitiateConsoleLogger()},
		FileLogger:    FileLogger{},
	}
}

func (z *ZeroLogger) WriteTrace(msg string) {
	f_logger, file := InitiateFileLogger()
	defer file.Close()
	z.FileLogger.File = file
	z.FileLogger.Logger = f_logger
	z.ConsoleLogger.Logger.Trace().Msg(msg)
	z.FileLogger.Logger.Trace().Msg(msg)
}

func (z *ZeroLogger) WriteError(msg string) {
	f_logger, file := InitiateFileLogger()
	defer file.Close()
	z.FileLogger.File = file
	z.FileLogger.Logger = f_logger
	z.ConsoleLogger.Logger.Error().Msg(msg)
	z.FileLogger.Logger.Error().Msg(msg)
}

func InitiateFileLogger() (zerolog.Logger, *os.File) {
	// opening log files
	file, err := os.OpenFile(
		"apilogs/myapp.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	return zerolog.New(file).With().Timestamp().Caller().Int("pid", os.Getpid()).Logger(), file
}

func InitiateConsoleLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr, TimeFormat: time.RFC3339,
	}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Logger()
}
