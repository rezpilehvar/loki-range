package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func CalculateTimeRange(timeRange string) (start string, end string, err error) {
	now := time.Now()
	switch {
	case timeRange == "today":
		{
			start = BeginningOfDay(now).Format(time.RFC3339)
			end = now.Format(time.RFC3339)
		}
	case timeRange == "yesterday":
		{
			yesterday := now.AddDate(0, 0, -1)
			start = BeginningOfDay(yesterday).Format(time.RFC3339)
			end = EndOfDay(yesterday).Format(time.RFC3339)
		}
	case strings.HasSuffix(timeRange, "d"):
		{
			daysStr, _ := strings.CutSuffix(timeRange, "d")
			days, err := strconv.Atoi(daysStr)
			if err != nil {
				return "", "", errors.New("invalid range format")
			}

			fromDate := now.AddDate(0, 0, -days)
			start = BeginningOfDay(fromDate).Format(time.RFC3339)
			end = now.Format(time.RFC3339)
		}
	case strings.HasSuffix(timeRange, "h"):
		{
			hoursStr, _ := strings.CutSuffix(timeRange, "h")
			hours, err := strconv.Atoi(hoursStr)
			if err != nil {
				return "", "", errors.New("invalid range format")
			}

			fromDate := now.Add(time.Duration(-hours) * time.Hour)
			start = fromDate.Format(time.RFC3339)
			end = now.Format(time.RFC3339)
		}
	case strings.HasSuffix(timeRange, "m"):
		{
			minutesStr, _ := strings.CutSuffix(timeRange, "m")
			minutes, err := strconv.Atoi(minutesStr)
			if err != nil {
				return "", "", errors.New("invalid range format")
			}

			fromDate := now.Add(time.Duration(-minutes) * time.Minute)
			start = fromDate.Format(time.RFC3339)
			end = now.Format(time.RFC3339)
		}
	default:
		{
			return "", "", errors.New("invalid range format")
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
