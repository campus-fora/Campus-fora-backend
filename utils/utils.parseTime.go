package utils

import "time"

func ParseTime(time_ string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.999999Z"
	paresedTime, err := time.Parse(layout,time_)
	if err != nil {
		return time.Time{}, err
	}
	return paresedTime, nil
}