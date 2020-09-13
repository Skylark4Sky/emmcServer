package utils

import (
	"time"
)

func GetTimestamp() int {
	return int(time.Now().Unix())
}

func GetTimestampMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetTimestampNs() int64 {
	return time.Now().UnixNano()
}

func TimeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}