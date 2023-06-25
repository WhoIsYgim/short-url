package utils

import "time"

const format = "2006-01-02 15:04:05"

func ExpireTimeString(expiration uint) string {
	return (time.Now().Add(time.Duration(expiration) * time.Hour)).Format(format)
}

func CurrentTimeString() string {
	return time.Now().Format(format)
}
