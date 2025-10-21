package main

import "log/slog"

type Logger interface {
	Info(message string, args ...any)
}

func New() Logger {
	return slog.New(nil)
}

func main() {
	f(12)
	f("")

	stdlog := slog.New(nil)
	stdlog.Info("hello", [2]any{"k", "v"}) // allocation

	cstlog := New()
	cstlog.Info("hello", "k", "v")
}

func f(x any) {
	_ = x
}
