package timesafer

import (
	"errors"
	"time"
)

type CET time.Location

func (c CET) Now() CETTime {
	loc := time.Location(c)
	return CETTime{time.Now().In(&loc)}
}

func (c CET) TimeAt(year int, month time.Month, day, hour, min, sec, nsec int) (CETTime, error) {
	loc := time.Location(c)
	t := time.Date(year, month, day, hour, min, sec, nsec, &loc)

	// go standard library will accept weird values and adjust the time returned
	// ex: time.Date(2022, 13, 32, 25, 0, 0, 0, 0, loc) will return 2023-02-02T01:00
	// we don't want this behavior because we think it's unexpected
	if year < 1 || t.Year() != year || t.Month() != month || t.Day() != day ||
		t.Hour() != hour || t.Minute() != min || t.Second() != sec || t.Nanosecond() != nsec {
		return CETTime{}, errors.New("time is invalid")
	}
	return CETTime{t}, nil
}

func NewCET() (CET, error) {
	loc, err := time.LoadLocation("Europe/Berlin")
	return CET(*loc), err
}

func MustCET() CET {
	cet, err := NewCET()
	if err != nil {
		panic(err)
	}
	return cet
}

type CETTime struct {
	t time.Time
}

func (c CETTime) RFC3339() string {
	return c.t.Format(time.RFC3339)
}

func (c CETTime) Date() CETDate {
	return CETDate{
		Year:  c.t.Year(),
		Month: c.t.Month(),
		Day:   c.t.Day(),
	}
}

type CETDate struct {
	Year  int
	Month time.Month
	Day   int
}
