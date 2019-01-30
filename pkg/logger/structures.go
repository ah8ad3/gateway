package logger

import "time"

type Log struct {
	Event string
	Description string

}

type UserLog struct {
	Log Log
	Ip string
	RequestUrl string
	Time time.Time
}

type SystemLog struct {
	Time time.Time
	Pkg string
	Log Log
}
