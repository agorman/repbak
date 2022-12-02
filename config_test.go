package repbak

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config, err := OpenConfig("./testdata/repbak.yaml")
	assert.Nil(t, err)

	assert.Equal(t, config.LogPath, "/var/log/repbak.log")
	assert.Equal(t, config.LogLevel, "error")

	assert.NotNil(t, config.HTTP)
	assert.Equal(t, config.HTTP.Addr, "0.0.0.0")
	assert.Equal(t, config.HTTP.Port, 4060)

	assert.NotNil(t, config.MySQL)
	assert.Equal(t, config.MySQL.Retention, 30)
	assert.Equal(t, config.MySQL.OutputPath, "/mnt/backups/mysql.dump")
	assert.Equal(t, config.MySQL.Schedule, "0 0 * * *")
	assert.Equal(t, config.MySQL.ExecutablePath, "mysqldump")
	assert.Equal(t, config.MySQL.ExecutableArgs, "--add-drop-database --all-databases -u user -ppass -h 127.0.0.1")
	assert.Equal(t, config.MySQL.TimeLimit, "8h")

	assert.NotNil(t, config.Email)
	assert.Equal(t, config.Email.Host, "mail.me.com")
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
}
