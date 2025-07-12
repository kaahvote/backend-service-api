package data

import (
	"errors"
	"time"
)

var ErrInvalidDateTimeFormat = errors.New("invalid date time format")

type DateTime string

func (dt DateTime) ToTime() (time.Time, error) {

	layout := "2006-01-02T15:04:00"
	t, err := time.Parse(layout, string(dt))

	if err != nil {
		return time.Time{}, ErrInvalidDateTimeFormat
	}

	return t, nil
}
