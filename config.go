package repbak

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

// Config is an object representation of the YAML configuration file.
type Config struct {
	// LogPath is the oath on disk where repbak log file. Defaults to /var/log/repbak.log.
	LogPath string `yaml:"log_path"`

	// LogLevel sets the level of logging. Valid levels are: panic, fatal, trace, debug, warn, info, and error. Defaults to error
	LogLevel string `yaml:"log_level"`

	// LibPath is the directory on disk where repbak lib files are stored. Defaults to /var/lib/repbak.
	LibPath string `yaml:"lib_path"`

	// The time format used when displaying backup stats. See formatting options in the go time.Time package.
	// Defaults to Mon Jan 02 03:04:05 PM MST
	TimeFormat string `yaml:"time_format"`

	// Retention is the number of logs and stats that are stored for each backup. If set to less than 0 no
	// logs or stats are saved. Defaults to 7.
	Retention int `yaml:"retention"`

	HTTP      *HTTP      `yaml:"http"`
	MySQLDump *MySQLDump `yaml:"mysqldump"`
	Email     *Email     `yaml:"email"`
}

// validate both validates the configuration and sets the default options.
func (c *Config) validate() error {
	if c.LogPath == "" {
		c.LogPath = "/var/log/repbak.log"
	}

	if c.LogLevel == "" {
		c.LogLevel = "error"
		log.SetLevel(log.ErrorLevel)
	} else {
		switch c.LogLevel {
		case "panic":
			log.SetLevel(log.PanicLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		case "trace":
			log.SetLevel(log.TraceLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		default:
			return fmt.Errorf("Invalid log_level: %s", c.LogLevel)
		}
	}

	if c.LibPath == "" {
		c.LibPath = "/var/lib/repbak"
	}

	if c.TimeFormat == "" {
		c.TimeFormat = "Mon Jan 02 03:04:05 PM MST"
	}

	if c.Retention == 0 {
		c.Retention = 7
	}

	if c.HTTP != nil {
		if c.HTTP.Addr == "" {
			c.HTTP.Addr = "127.0.0.1"
		}

		if c.HTTP.Port == 0 {
			c.HTTP.Port = 4060
		}
	}

	if c.Email == nil {
		return errors.New("Missing required email configuration")
	}

	if c.Email.Host == "" {
		return errors.New("Missing required host entry for email")
	}

	if c.Email.Port == 0 {
		c.Email.Port = 25
	}

	// StartTLS takes presidence over SSL
	if c.Email.StartTLS {
		c.Email.SSL = false
	}

	if c.Email.Subject == "" {
		c.Email.Subject = "Database Replication Failure"
	}

	if c.Email.From == "" {
		return errors.New("Missing required from entry for email")
	}

	if len(c.Email.To) == 0 {
		return errors.New("Missing required to entry for email")
	}

	if c.Email.HistorySubject == "" {
		c.Email.HistorySubject = "Database Backup History"
	}

	if c.MySQLDump == nil {
		return errors.New("Missing required mysqldump configuration")
	}

	if c.MySQLDump.Retention == 0 {
		c.MySQLDump.Retention = 7
	}

	if c.MySQLDump.OutputPath == "" {
		return errors.New("Missing required output_path entry for mysqldump")
	}

	if c.MySQLDump.Schedule == "" {
		return errors.New("Missing required schedule entry for mysqldump")
	}

	if c.MySQLDump.ExecutablePath == "" {
		c.MySQLDump.ExecutablePath = "mysqldump"
	}

	if c.MySQLDump.ExecutableArgs == "" {
		c.MySQLDump.ExecutableArgs = "--add-drop-database --all-databases"
	}

	if c.MySQLDump.TimeLimit != "" {
		var err error
		c.MySQLDump.timeLimit, err = time.ParseDuration((c.MySQLDump.TimeLimit))
		if err != nil {
			return fmt.Errorf("Failed to parse mysqldump time_limit: %w", err)
		}
	}

	return nil

}

// MySQL defines how a mysql backup will be created.
type MySQLDump struct {
	// Retention is the number of backups to keep before rotating old backups out. Defaults to 7.
	Retention int `yaml:"retention"`

	// OutputPath is the path where backups will be stored.
	OutputPath string `yaml:"output_path"`

	// Schedule is the cron expression that defines when backups are created.
	Schedule string `yaml:"schedule"`

	// ExecutablePath is the path to the tool used to create the mysql backup. Defaults to mysqldump.
	ExecutablePath string `yaml:"executable_path"`

	// ExecutableArgs are the arguments passed to the executable used to create the mysql backup. Defaults to --add-drop-database --all-databases.
	ExecutableArgs string `yaml:"executable_args"`

	// TimeLimit is an optional limit to the time it takes to run the backup.
	TimeLimit string `yaml:"time_limit"`
	timeLimit time.Duration
}

// HTTP defines the configuration for http health checks.
type HTTP struct {
	// The address the http server will listen on.
	Addr string `yaml:"addr"`

	// The port the http server will listen on.
	Port int `yaml:"port"`
}

type Email struct {
	// Host is the hostname or IP of the SMTP server.
	Host string `yaml:"host"`

	// Port is the port of the SMTP server.
	Port int `yaml:"port"`

	// User is the username used to authenticate.
	User string `yaml:"user"`

	// Pass is the password used to authenticate.
	Pass string `yaml:"pass"`

	// StartTLS enables TLS security. If both StartTLS and SSL are true then StartTLS will be used.
	StartTLS bool `yaml:"starttls"`

	// Skip verifying the server's certificate chain and host name.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`

	// SSL enables SSL security. If both StartTLS and SSL are true then StartTLS will be used.
	SSL bool `yaml:"ssl"`

	// Optional subject field for notification emails
	Subject string `yaml:"subject"`

	// From is the email address the email will be sent from.
	From string `yaml:"from"`

	// To is an array of email addresses for which emails will be sent.
	To []string `yaml:"to"`

	// HistorySubject is an optional subject to use when sending sync history emails.
	HistorySubject string `yaml:"history_subject"`

	// HistorySchedule is a cron expression. If set then an email with sync history will be sent based on the schedule.
	HistorySchedule string `yaml:"history_schedule"`

	// HistoryTemplate is an optional path to an email template to use when sending history emails. If not set uses the default template.
	HistoryTemplate string `yaml:"history_template"`

	// OnFailure will send an email for each backup failure if true.
	OnFailure bool `yaml:"on_failure"`
}

// OpenConfig returns a new Config option by reading the YAML file at path. If the file
// doesn't exist, can't be read, is invalid YAML, or doesn't match the repbak spec then
// an error is returned.
func OpenConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	return config, config.validate()
}
