package repbak

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config, err := OpenConfig("./testdata/repbak.yaml")
	assert.Nil(t, err)

	assert.Equal(t, config.LogPath, "/tmp/repbak.log")
	assert.Equal(t, config.LogLevel, "fatal")
	assert.Equal(t, config.LibPath, "/tmp/repbak")
	assert.Equal(t, config.Retention, 10)
	assert.Equal(t, config.TimeFormat, "2006-01-02T15:04:05Z07:00")

	assert.NotNil(t, config.HTTP)
	assert.Equal(t, config.HTTP.Addr, "0.0.0.0")
	assert.Equal(t, config.HTTP.Port, 4060)

	assert.NotNil(t, config.MySQLDump)
	assert.Equal(t, config.MySQLDump.Retention, 30)
	assert.Equal(t, config.MySQLDump.OutputPath, "/mnt/backups/mysql.dump")
	assert.Equal(t, config.MySQLDump.Schedule, "0 0 * * *")
	assert.Equal(t, config.MySQLDump.ExecutablePath, "mysqldump")
	assert.Equal(t, config.MySQLDump.ExecutableArgs, "--add-drop-database --all-databases -u user -ppass -h 127.0.0.1")
	assert.Equal(t, config.MySQLDump.TimeLimit, "8h")

	assert.NotNil(t, config.Email)
	assert.Equal(t, config.Email.Host, "1.1.1.1.1")
	assert.Equal(t, config.Email.Port, 587)
	assert.Equal(t, config.Email.User, "me")
	assert.Equal(t, config.Email.Pass, "pass")
	assert.Equal(t, config.Email.StartTLS, true)
	assert.Equal(t, config.Email.InsecureSkipVerify, false)
	assert.Equal(t, config.Email.SSL, false)
	assert.Equal(t, config.Email.From, "me@me.com")
	assert.Contains(t, config.Email.To, "they@me.com")
	assert.Contains(t, config.Email.To, "them@me.com")
	assert.Equal(t, config.Email.Subject, "Database Backup Failure")
	assert.Equal(t, config.Email.OnFailure, true)
}

func TestConfigDefaults(t *testing.T) {
	config, err := OpenConfig("./testdata/defaults.yaml")
	assert.Nil(t, err)

	config.HTTP = &HTTP{}

	err = config.validate()
	assert.Nil(t, err)
	assert.Equal(t, config.LogPath, "/var/log/repbak.log")
	assert.Equal(t, config.LogLevel, "error")
	assert.Equal(t, config.LibPath, "/var/lib/repbak")
	assert.Equal(t, config.Retention, 7)
	assert.Equal(t, config.TimeFormat, "Mon Jan 02 03:04:05 PM MST")

	assert.Equal(t, config.Email.Port, 25)
	assert.Equal(t, config.Email.StartTLS, false)
	assert.Equal(t, config.Email.SSL, false)
	assert.Equal(t, config.Email.Subject, "Database Replication Failure")
	assert.Equal(t, config.Email.OnFailure, false)

	assert.Equal(t, config.MySQLDump.Retention, 7)
	assert.Equal(t, config.MySQLDump.ExecutablePath, "mysqldump")
	assert.Equal(t, config.MySQLDump.ExecutableArgs, "--add-drop-database --all-databases")

	assert.Equal(t, config.HTTP.Addr, "127.0.0.1")
	assert.Equal(t, config.HTTP.Port, 4060)
}

func TestConfigRequired(t *testing.T) {
	config := &Config{}

	err := config.validate()
	assert.Error(t, err)

	config.Email = &Email{}
	err = config.validate()
	assert.Error(t, err)

	config.Email.Host = "1.1.1.1.1"
	err = config.validate()
	assert.Error(t, err)

	config.Email.From = "me@me.com"
	err = config.validate()
	assert.Error(t, err)

	config.Email.To = []string{"you@me.com"}
	err = config.validate()
	assert.Error(t, err)

	config.MySQLDump = &MySQLDump{}
	err = config.validate()
	assert.Error(t, err)

	config.MySQLDump.OutputPath = "/tmp/mysql.dump"
	err = config.validate()
	assert.Error(t, err)

	config.MySQLDump.Schedule = "0 0 * * *"
	err = config.validate()
	assert.Nil(t, err)
}

func TestConfigBadPath(t *testing.T) {
	_, err := OpenConfig("./testdata/notexist.yaml")
	assert.Error(t, err)
}

func TestConfigLogLevel(t *testing.T) {
	config, err := OpenConfig("./testdata/defaults.yaml")
	assert.Nil(t, err)

	config.LogLevel = "panic"
	err = config.validate()
	assert.Nil(t, err)

	config.LogLevel = "fatal"
	err = config.validate()
	assert.Nil(t, err)

	config.LogLevel = "trace"
	err = config.validate()
	assert.Nil(t, err)

	config.LogLevel = "debug"
	err = config.validate()
	assert.Nil(t, err)

	config.LogLevel = "warn"
	err = config.validate()
	assert.Nil(t, err)

	config.LogLevel = "info"
	err = config.validate()
	assert.Nil(t, err)

	config.LogLevel = "error"
	err = config.validate()
	assert.Nil(t, err)

	config.LogLevel = "bad"
	err = config.validate()
	assert.Error(t, err)
}

func TestConfigOpenInvalid(t *testing.T) {
	_, err := OpenConfig("./testdata/invalid.yaml")
	assert.Error(t, err)
}

func TestConfigBadTimeLimit(t *testing.T) {
	config, err := OpenConfig("./testdata/repbak.yaml")
	assert.Nil(t, err)

	config.MySQLDump.TimeLimit = "asdf"
	err = config.validate()
	assert.Error(t, err)
}
