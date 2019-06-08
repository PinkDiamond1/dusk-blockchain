package logger

import (
	log "github.com/sirupsen/logrus"
)

// This package contains a logrus hook to capture all errors and warnings so they can be sent also to the log file read by the monitor without the need to discriminate by lines
func NotifyErrors(proc *LogProcessor) {
	log.AddHook(proc)
}
