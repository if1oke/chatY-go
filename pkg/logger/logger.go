package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type ILogger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
}

type logrusAdapter struct{ *logrus.Logger }

func (l *logrusAdapter) Info(args ...any)                  { l.Logger.Info(args...) }
func (l *logrusAdapter) Infof(format string, args ...any)  { l.Logger.Infof(format, args...) }
func (l *logrusAdapter) Error(args ...any)                 { l.Logger.Error(args...) }
func (l *logrusAdapter) Errorf(format string, args ...any) { l.Logger.Errorf(format, args...) }
func (l *logrusAdapter) Warn(args ...any)                  { l.Logger.Warn(args...) }
func (l *logrusAdapter) Warnf(format string, args ...any)  { l.Logger.Warnf(format, args...) }
func (l *logrusAdapter) Debug(args ...any)                 { l.Logger.Debug(args...) }
func (l *logrusAdapter) Debugf(format string, args ...any) { l.Logger.Debugf(format, args...) }

func NewLogger() ILogger {
	log := logrus.New()

	// Формат
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	// Дефолт уровень
	log.SetLevel(logrus.InfoLevel)

	// Лог в файл
	file, err := os.OpenFile("logs/chat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Warn("Could not log to file, using default stderr")
	} else {
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	}

	return &logrusAdapter{log}
}
