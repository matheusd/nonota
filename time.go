package nonota

import (
	"time"
)

// StartOfBilling returns the time of the start of billing for a given time.
func StartOfBilling(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
}

func EndOfBilling(t time.Time) time.Time {
	return EndOfDay(time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.Local))
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
}

func StartOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	oneday, _ := time.ParseDuration("24h")
	startDay := t.Add(oneday * time.Duration(weekday) * -1)
	return StartOfDay(startDay)
}

func EndOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	oneday, _ := time.ParseDuration("24h")
	endDay := t.Add(oneday * (6 - time.Duration(weekday)))
	return EndOfDay(endDay)
}