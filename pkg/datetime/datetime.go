package datetime

import "time"

const TimeFormat = "2006-01-02 15:04:05"

func NowLocalString() string {
	currentTime := time.Now().Local()
	layoutISO := currentTime.Format(TimeFormat)
	return layoutISO
}

func Parse(s string) (time.Time, error) {
	return time.Parse(TimeFormat, s)
}
