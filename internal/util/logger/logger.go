package logger

import (
	"log/slog"
)

const (
	levelDebug   = "debug"
	levelInfo    = "info"
	levelWarning = "warn"
	levelError   = "error"

	defaultLevel = slog.LevelWarn
)

type Logger struct {
	log *slog.Logger
}

//func New(cfg config.LoggerConfig) *Logger {
//	var output = os.Stdout
//
//	if cfg.LogFile != nil {
//		tempOutput, err := os.OpenFile(*cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
//		if err == nil {
//			output = tempOutput
//		}
//	}
//
//	opts := &slog.HandlerOptions{
//		Level: getLevel(cfg.Level),
//	}
//
//	return &Logger{
//		log: slog.New(slog.NewJSONHandler(output, opts)),
//	}
//}

func getLevel(lvl string) slog.Level {
	var level slog.Level

	switch lvl {
	case levelDebug:
		level = slog.LevelDebug
	case levelInfo:
		level = slog.LevelInfo
	case levelWarning:
		level = slog.LevelWarn
	case levelError:
		level = slog.LevelError
	default:
		level = defaultLevel
	}

	return level
}
