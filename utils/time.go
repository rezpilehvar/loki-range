package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func Format(dateTime *time.Time) string {
	return dateTime.Format(time.RFC3339)
}
func CalculateTimeRange(timeRange string) (start *time.Time, end *time.Time, err error) {
	now := time.Now()
	switch {
	case timeRange == "today":
		{
			startTime := BeginningOfDay(now)
			start = &startTime
			end = &now
		}
	case timeRange == "yesterday":
		{
			yesterday := now.AddDate(0, 0, -1)
			startTime := BeginningOfDay(yesterday)
			endTime := EndOfDay(yesterday)
			start = &startTime
			end = &endTime
		}
	case strings.HasSuffix(timeRange, "d"):
		{
			daysStr, _ := strings.CutSuffix(timeRange, "d")
			days, err := strconv.Atoi(daysStr)
			if err != nil {
				return nil, nil, errors.New("invalid range format")
			}

			fromDate := now.AddDate(0, 0, -days)
			startTime := BeginningOfDay(fromDate)

			start = &startTime
			end = &now
		}
	case strings.HasSuffix(timeRange, "h"):
		{
			hoursStr, _ := strings.CutSuffix(timeRange, "h")
			hours, err := strconv.Atoi(hoursStr)
			if err != nil {
				return nil, nil, errors.New("invalid range format")
			}

			fromDate := now.Add(time.Duration(-hours) * time.Hour)
			start = &fromDate
			end = &now
		}
	case strings.HasSuffix(timeRange, "m"):
		{
			minutesStr, _ := strings.CutSuffix(timeRange, "m")
			minutes, err := strconv.Atoi(minutesStr)
			if err != nil {
				return nil, nil, errors.New("invalid range format")
			}

			fromDate := now.Add(time.Duration(-minutes) * time.Minute)
			start = &fromDate
			end = &now
		}
	case strings.HasSuffix(timeRange, "s"):
		{
			secondsStr, _ := strings.CutSuffix(timeRange, "s")
			seconds, err := strconv.Atoi(secondsStr)
			if err != nil {
				return nil, nil, errors.New("invalid range format")
			}

			fromDate := now.Add(time.Duration(-seconds) * time.Second)
			start = &fromDate
			end = &now
		}
	default:
		{
			return nil, nil, errors.New("invalid range format")
		}
	}

	return start, end, nil
}

func BeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

func BeginningYesterday(t time.Time) time.Time {
	return BeginningOfDay(t)
}
