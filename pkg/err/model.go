package err

import (
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
)

// Err custom error structure for app
type Err struct {
	Message  string
	Critical bool
}

// Log for log the exception
// category for system or user log or exception
func (err Err) Log(category string) Err {
	if category == "system" {
		logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
			Description: err.Message}})
	} else if category == "user" {
		logger.SetSysLog(logger.SystemLog{Pkg: "db", Time: time.Now(), Log: logger.Log{Event: "critical",
			Description: err.Message}})
	}
	return err
}
