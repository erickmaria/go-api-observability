package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func NewLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
