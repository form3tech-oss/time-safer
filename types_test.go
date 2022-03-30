package timesafer_test

import (
	"math"
	"testing"
	"time"

	"github.com/form3tech-oss/time-safer/pkg/timesafer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCET_TimeIsInLocation(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	require.NoError(t, err, "failed to load timezone location")
	now := time.Now().In(loc)

	cet, err := timesafer.NewCET()
	require.NoError(t, err)
	saferNow := cet.Now()
	assert.Equal(t, now.Format(time.RFC3339), saferNow.RFC3339())
}

func TestCET_TimeAtValidTime(t *testing.T) {
	cet := timesafer.MustCET()
	cetTime, err := cet.TimeAt(2022, time.March, 2, 15, 33, 40, 0)
	assert.Equal(t, cet, cetTime.CET())
	require.NoError(t, err)
	assert.Equal(t, "2022-03-02T15:33:40+01:00", cetTime.RFC3339())
}

func TestCET_TimeAtInvalidTime(t *testing.T) {
	cet := timesafer.MustCET()

	t.Run("invalid year", func(t *testing.T) {
		_, err := cet.TimeAt(0, 11, 2, 23, 59, 59, 0)
		require.Error(t, err)
	})

	t.Run("invalid month", func(t *testing.T) {
		_, err := cet.TimeAt(2022, 13, 2, 23, 59, 59, 0)
		require.Error(t, err)
		_, err = cet.TimeAt(2022, -1, 2, 23, 59, 59, 0)
		require.Error(t, err)
	})

	t.Run("invalid day", func(t *testing.T) {
		_, err := cet.TimeAt(2022, 11, 32, 23, 59, 59, 0)
		require.Error(t, err)
		_, err = cet.TimeAt(2022, 11, -1, 23, 59, 59, 0)
		require.Error(t, err)
	})

	t.Run("invalid hour", func(t *testing.T) {
		_, err := cet.TimeAt(2022, 11, 2, 25, 59, 59, 0)
		require.Error(t, err)
		_, err = cet.TimeAt(2022, 11, 2, -1, 59, 59, 0)
		require.Error(t, err)
	})

	t.Run("invalid minute", func(t *testing.T) {
		_, err := cet.TimeAt(2022, 11, 2, 23, 69, 59, 0)
		require.Error(t, err)
		_, err = cet.TimeAt(2022, 11, 2, 23, -1, 59, 0)
		require.Error(t, err)
	})

	t.Run("invalid second", func(t *testing.T) {
		_, err := cet.TimeAt(2022, 11, 2, 23, 59, 69, 0)
		require.Error(t, err)
		_, err = cet.TimeAt(2022, 11, 2, 23, 59, -1, 0)
		require.Error(t, err)
	})

	t.Run("invalid nanosecond", func(t *testing.T) {
		_, err := cet.TimeAt(2022, 11, 2, 23, 59, 59, math.MaxInt)
		require.Error(t, err)
		_, err = cet.TimeAt(2022, 11, 2, 23, 59, 59, -1)
		require.Error(t, err)
	})
}

func TestCET_TimeToDate(t *testing.T) {
	cet := timesafer.MustCET()

	t.Run("middle of the day", func(t *testing.T) {
		cetTime, err := cet.TimeAt(2022, time.March, 2, 12, 33, 40, 0)
		require.NoError(t, err)
		cetDate := cetTime.Date()
		assert.Equal(t, 2022, cetDate.Year)
		assert.Equal(t, time.March, cetDate.Month)
		assert.Equal(t, 2, cetDate.Day)
	})

	t.Run("just after midnight", func(t *testing.T) {
		cetTime, err := cet.TimeAt(2022, time.March, 2, 0, 1, 0, 0)
		require.NoError(t, err)
		cetDate := cetTime.Date()
		assert.Equal(t, 2022, cetDate.Year)
		assert.Equal(t, time.March, cetDate.Month)
		assert.Equal(t, 2, cetDate.Day)
	})

	t.Run("just before midnight", func(t *testing.T) {
		cetTime, err := cet.TimeAt(2022, time.March, 2, 23, 59, 59, 0)
		require.NoError(t, err)
		cetDate := cetTime.Date()
		assert.Equal(t, 2022, cetDate.Year)
		assert.Equal(t, time.March, cetDate.Month)
		assert.Equal(t, 2, cetDate.Day)
	})
}

func TestCET_DateMarshalText(t *testing.T) {
	cet := timesafer.MustCET()
	date, err := cet.DateAt(2022, 1, 1)
	require.NoError(t, err)
	expected := "2022-01-01"
	actual, err := date.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestCET_DateUnmarshalTextSuccess(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		date := "2022-03-16"
		cetDate := timesafer.CETDate{}
		err := cetDate.UnmarshalText([]byte(date))
		require.NoError(t, err)
		assert.Equal(t, 2022, cetDate.Year)
		assert.Equal(t, time.March, cetDate.Month)
		assert.Equal(t, 16, cetDate.Day)
	})
	t.Run("valid date leap year", func(t *testing.T) {
		date := "2020-02-29"
		cetDate := timesafer.CETDate{}
		err := cetDate.UnmarshalText([]byte(date))
		require.NoError(t, err)
		assert.Equal(t, 2020, cetDate.Year)
		assert.Equal(t, time.February, cetDate.Month)
		assert.Equal(t, 29, cetDate.Day)
	})
}

func TestCET_DateUnmarshalTextFailure(t *testing.T) {
	tests := []struct {
		name, date string
	}{
		{
			name: "long february",
			date: "2020-02-30",
		},
		{
			name: "zero year",
			date: "0-01-01",
		},
		{
			name: "zero month",
			date: "2020-00-01",
		},
		{
			name: "zero day",
			date: "2020-01-00",
		},
		{
			name: "invalid format no hyphens",
			date: "20200101",
		},
		{
			name: "invalid format no leading zeros",
			date: "2020-1-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cetDate := timesafer.CETDate{}
			err := cetDate.UnmarshalText([]byte(tt.date))
			require.Error(t, err)
		})
	}
}

func TestCET_TimeMarshalText(t *testing.T) {
	cet := timesafer.MustCET()
	tm, err := cet.TimeAt(2022, 1, 1, 1, 1, 1, 123000000)
	require.NoError(t, err)
	actual, err := tm.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2022-01-01T01:01:01.123+01:00", string(actual))
}

func TestCET_TimeUnmarshalTextSuccess(t *testing.T) {
	tm := timesafer.CETTime{}
	require.NoError(t, tm.UnmarshalText([]byte("2022-01-01T01:01:01.123+01:00")))
	assert.Equal(t, 2022, tm.Year())
	assert.Equal(t, time.January, tm.Month())
	assert.Equal(t, 1, tm.Day())
	assert.Equal(t, 1, tm.Hour())
	assert.Equal(t, 1, tm.Minute())
	assert.Equal(t, 1, tm.Second())
	assert.Equal(t, 123000000, tm.Nanosecond())
}

func TestCET_TimeUnmarshalTextSuccessDifferentTimezone(t *testing.T) {
	tm := timesafer.CETTime{}
	require.NoError(t, tm.UnmarshalText([]byte("2022-12-31T23:59:00Z")))
	assert.Equal(t, 2023, tm.Year())
	assert.Equal(t, time.January, tm.Month())
	assert.Equal(t, 1, tm.Day())
	assert.Equal(t, 0, tm.Hour())
	assert.Equal(t, 59, tm.Minute())
	assert.Equal(t, 0, tm.Second())
	assert.Equal(t, 0, tm.Nanosecond())
}

func TestCET_TimeUnmarshalTextSuccessMissingTimezone(t *testing.T) {
	tm := timesafer.CETTime{}
	assert.NoError(t, tm.UnmarshalText([]byte("0001-12-31T23:59:00")))
	assert.Equal(t, 1, tm.Year())
	assert.Equal(t, time.December, tm.Month())
	assert.Equal(t, 31, tm.Day())
	assert.Equal(t, 23, tm.Hour())
	assert.Equal(t, 59, tm.Minute())
	assert.Equal(t, 0, tm.Second())
	assert.Equal(t, 0, tm.Nanosecond())
}

func TestCET_TimeUnmarshalTextFailure(t *testing.T) {
	tm := timesafer.CETTime{}
	assert.Error(t, tm.UnmarshalText([]byte("2022-12-31T23:59")))
	assert.Error(t, tm.UnmarshalText([]byte("2022-12-31T23:59.123")))
	assert.Error(t, tm.UnmarshalText([]byte("1-12-31T23:59:00")))
}
