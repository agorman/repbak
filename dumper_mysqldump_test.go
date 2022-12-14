package repbak

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMySQLDumpDumper(t *testing.T) {
	config, err := OpenConfig("./testdata/repbak.yaml")
	assert.Nil(t, err)

	dumper := NewMySQLDumpDumper(config)

	stat := dumper.Dump()
	assert.Error(t, stat.Error)

	dumper.Stop()
}
