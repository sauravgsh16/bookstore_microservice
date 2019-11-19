package dates

import (
	"time"
)

const (
	apiDateLayout = "2006-01-20T15:04:05Z"
)

// GetNow returns current time in UTC
func GetNow() time.Time {
	return time.Now().UTC()
}

// GetNowString returns current time in string
func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}
