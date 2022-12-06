package repbak

import (
	"errors"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestEmailNotifier(t *testing.T) {
	config, err := OpenConfig("./testdata/repbak.yaml")
	assert.Nil(t, err)

	notifier := NewEmailNotifier(config)

	stat := NewStat("TEST", "Mon Jan 02 03:04:05 PM MST").Finish(nil)

	err = notifier.Notify(stat)
	assert.Error(t, err)

	stat = stat.Finish(errors.New("ERROR"))

	err = notifier.Notify(stat)
	assert.Error(t, err)
}
