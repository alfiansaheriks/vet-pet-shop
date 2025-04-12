package utils

import "time"

func ParseTimeStampToDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func ParseDatetoTimestamp(date string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
