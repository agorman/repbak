package repbak

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestRepBak(t *testing.T) {
	config, err := OpenConfig("./testdata/repbak.yaml")
	assert.Nil(t, err)

	db, err := NewBoltDB(config)
	assert.Nil(t, err)
	defer db.Close()

	notifier := NewEmailNotifier(config)

	dumper := NewMySQLDumpDumper(config)

	rm := New(config, db, dumper, notifier)
	rm.Start()
	rm.Start()
	rm.Stop()
	rm.Stop()
}
