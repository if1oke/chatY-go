package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewClientLogger() ILogger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel)
	file, err := os.OpenFile("logs/chat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Warn("Could not log to file, skip...")
	}
	log.SetOutput(file)
	return &logrusAdapter{log}
}
