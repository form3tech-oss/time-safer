package timesafer_test

import (
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
