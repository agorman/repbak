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

	err = notifier.Notify(nil)
	assert.Error(t, err)

	err = notifier.Notify(errors.New("ERROR"))
	assert.Error(t, err)
}
