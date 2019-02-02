package logger

import "time"

// Log structure of the system base struct
type Log struct {
	Event       string
	Description string
}

// UserLog that extend Log and for the api logs
type UserLog struct {
	Log        Log
	IP         string
	RequestURL string
	Time       time.Time
}

// SystemLog that extend from base Log and for System failure and logs
type SystemLog struct {
	Time time.Time
	Pkg  string
	Log  Log
}
