package utils

import "time"

var appStartTime = time.Now()

func SetAppStartTime(t time.Time) {
	appStartTime = t
}

func GetAppStartTime() time.Time {
	return appStartTime
}

func GetUptime() time.Duration {
	return time.Since(appStartTime)
}

func GetUptimeString() string {
	return GetUptime().String()
}
