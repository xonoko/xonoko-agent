package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"os"
)

var logger *logrus.Entry

func init() {

	logLevel := logrus.DebugLevel
	envLogLevel := os.Getenv("XENOKO_LOG_LEVEL")
	switch envLogLevel {
	case "info":
		logLevel = logrus.InfoLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "error":
		logLevel = logrus.ErrorLevel
	}
	logrus.SetLevel(logLevel)
	initLogger := logrus.New()

	if os.Getenv("XENOKO_LOG_SLACK") == "true" {
		slackHook := &slackrus.SlackrusHook{
			HookURL : os.Getenv("XENOKO_SLACK_LOG_HOOK"),
			AcceptedLevels: slackrus.LevelThreshold(logLevel),
			Channel:        os.Getenv("XENOKO_SLACK_AGENT_CHANNEL"),
			IconEmoji:      ":ghost:",
			Username:       "logbot",
		}
		initLogger.Hooks.Add(slackHook)
	}
	logger = initLogger.WithFields(logrus.Fields{"program":"agent-dev"})
}
