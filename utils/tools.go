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
