package logger

import (
	"github.com/deqode/dq-vault/config"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Log - for logging of messages
func Log(logger log.Logger, level string, messages ...string) {
	message := strings.Join(messages, " ")

	switch level {
	case config.Debug:
		logger.Debug("\n"+timestamp(), "["+level+" ] ", message)
	case config.Info:
		logger.Info("\n"+timestamp(), "["+level+" ] ", message)
	case config.Error:
		logger.Error("\n"+timestamp(), "["+level+" ] ", message)
	case config.Fatal:
		logger.Fatal("\n"+timestamp(), "["+level+" ] ", message)
	default:
		return
	}
}

// Timestamp - to identify time of occurence of an event
// returns current timestamp
// example - 2018-08-07 12:04:46.456601867 +0000 UTC m=+0.000753626
func timestamp() string {
	return time.Now().Format(time.RFC3339)
}
