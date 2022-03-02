package timesafer

import "time"

type CET time.Location

func (c CET) Now() CETTime {
	loc := time.Location(c)
	return CETTime{time.Now().In(&loc)}
}

type CETTime struct {
	t time.Time
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

func (c CETTime) RFC3339() string {
	return c.t.Format(time.RFC3339)
}
