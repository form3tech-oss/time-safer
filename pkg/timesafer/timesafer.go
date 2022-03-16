package timesafer

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
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

func (c CET) DateAt(year int, month time.Month, day int) (CETDate, error) {
	t, err := c.TimeAt(year, month, day, 0, 0, 0, 0)
	return t.Date(), err
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

func (d *CETDate) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)), nil
}

func (d *CETDate) UnmarshalText(text []byte) error {
	chunks := bytes.Split(text, []byte("-"))
	if len(chunks) != 3 {
		return fmt.Errorf("invalid date")
	}
	if len(chunks[1]) != 2 || len(chunks[2]) != 2 {
		return fmt.Errorf("invalid date")
	}
	var err error
	d.Year, err = validateYear(string(chunks[0]))
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	d.Month, err = validateMonth(string(chunks[1]))
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	d.Day, err = validateDay(string(chunks[2]), d.Month, d.Year)
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}
	return nil
}

func validateYear(year string) (int, error) {
	var err error
	y, err := strconv.Atoi(year)
	if err != nil {
		return 0, fmt.Errorf("invalid year: %w", err)
	}
	if y < 1 {
		return 0, fmt.Errorf("year should be greater than 0")
	}
	return y, nil
}

func validateMonth(month string) (time.Month, error) {
	m, err := strconv.Atoi(month)
	if err != nil {
		return 0, fmt.Errorf("invalid month: %w", err)
	}
	if m > int(time.December) {
		return 0, fmt.Errorf("month should be less than 13")
	}
	if m < int(time.January) {
		return 0, fmt.Errorf("month should be greater than 0")
	}
	return time.Month(m), nil
}

func validateDay(day string, month time.Month, year int) (int, error) {
	d, err := strconv.Atoi(day)
	if err != nil {
		return 0, fmt.Errorf("invalid day: %w", err)
	}
	if d < 1 {
		return 0, fmt.Errorf("day should be greater than 0")
	}
	daysInMonth := daysIn(month, year)
	if d > daysInMonth {
		return 0, fmt.Errorf("there is only %d days in %s", daysInMonth, month)
	}
	return d, nil
}

// TODO: this is taken from standard library time
// TODO: handle licensing
func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// daysBefore[m] counts the number of days in a non-leap year
// before month m begins. There is an entry for m=12, counting
// the number of days before January of next year (365).
var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

func daysIn(m time.Month, year int) int {
	if m == time.February && isLeap(year) {
		return 29
	}
	return int(daysBefore[m] - daysBefore[m-1])
}
