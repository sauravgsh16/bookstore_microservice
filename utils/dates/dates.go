package dates

import (
	"fmt"
	"time"
)

const (
	apiDateLayout   = "2006-01-20T15:04:05Z"
	apiDateDBLayout = "2006-01-2 15:04:05"
)

// GetNow returns current time in UTC
func GetNow() time.Time {
	return time.Now().UTC()
}

// GetNowString returns current time in string
func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

// GetNowDBString returns current time db DATETIME format
func GetNowDBString() string {
	fmt.Printf("DATE: %s\n", GetNow().Format(apiDateDBLayout))
	return GetNow().Format(apiDateDBLayout)
}
