package repbak

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStat(t *testing.T) {
	stat := NewStat("SUCCESS", "Mon Jan 02 03:04:05 PM MST")
	assert.Equal(t, stat.Name, "SUCCESS")
	assert.False(t, stat.Success)
	assert.Equal(t, stat.End, "")
	assert.Equal(t, stat.Duration, time.Duration(0))
	assert.Equal(t, stat.Skip, false)

	stat = stat.Finish(nil)
	assert.Equal(t, stat.Name, "SUCCESS")
	assert.True(t, stat.Success)
	assert.NotEqual(t, stat.End, "")
	assert.NotEqual(t, stat.Duration, time.Duration(0))
	assert.Equal(t, stat.Skip, false)
}

func TestErrorStat(t *testing.T) {
	stat := NewStat("FAIL", "Mon Jan 02 03:04:05 PM MST")
	assert.Equal(t, stat.Name, "FAIL")
	assert.False(t, stat.Success)
	assert.Equal(t, stat.End, "")
	assert.Equal(t, stat.Duration, time.Duration(0))
	assert.Equal(t, stat.Skip, false)

	err := errors.New("fail")
	stat = stat.Finish(err)
	assert.Equal(t, stat.Name, "FAIL")
	assert.False(t, stat.Success)
	assert.NotEqual(t, stat.End, "")
	assert.NotEqual(t, stat.Duration, time.Duration(0))
	assert.Equal(t, stat.Skip, false)
	assert.Equal(t, stat.Error, err)
}
