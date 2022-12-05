package repbak

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestRepBak(t *testing.T) {
	config, err := OpenConfig("./testdata/repbak.yaml")
	assert.Nil(t, err)

	dumper := NewMySQLDumpDumper(config)

	notifier := NewEmailNotifier(config)

	rm := New(config, dumper, notifier)
	rm.Start()
	rm.Start()
	rm.Stop()
	rm.Stop()
}
