package main

import (
	"log"
	"os"
)

type LogLevel int64

const (
	DEBUG LogLevel = iota
	WARNING
	INFO
	ERROR
)

type Logger struct {
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
	Level         LogLevel
}

func (l *Logger) init(args ...LogLevel) {
	if len(args) > 0 {
		l.Level = args[0]
	} else {
		l.Level = INFO
	}

	l.InfoLogger = log.New(os.Stdout, "[ INFO ]: ", log.LstdFlags|log.Lmsgprefix)
	l.ErrorLogger = log.New(os.Stdout, "[ ERROR ]: ", log.LstdFlags|log.Lmsgprefix)
	l.WarningLogger = log.New(os.Stdout, "[ WARNING ]: ", log.LstdFlags|log.Lmsgprefix)
	l.DebugLogger = log.New(os.Stdout, "[ DEBUG ]: ", log.LstdFlags|log.Lmsgprefix)
}

func (l Logger) Debug(msg string) {
	if l.Level == DEBUG {
		l.DebugLogger.Println(msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.Level <= WARNING {
		l.WarningLogger.Println(msg)
	}
}

func (l Logger) Info(msg string) {
	if l.Level <= INFO {
		l.InfoLogger.Println(msg)
	}
}

func (l Logger) Error(msg string) {
	l.ErrorLogger.Println(msg)
}
