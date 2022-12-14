[![Build Status](https://github.com/agorman/repbak/workflows/repbak/badge.svg)](https://github.com/agorman/repbak/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/agorman/repbak)](https://goreportcard.com/report/github.com/agorman/repbak)
[![GoDoc](https://godoc.org/github.com/agorman/repbak?status.svg)](https://godoc.org/github.com/agorman/repbak)
[![codecov](https://codecov.io/gh/agorman/repbak/branch/main/graph/badge.svg)](https://codecov.io/gh/agorman/repbak)

# Repbak


Repbak is a simple database backup tool made specifically to backup replicated databases. Repbak will send notifications based on database backup failure. Repbak also optionally support HTTP healthchecks for liveness.


# Supported Dumpers


- mysqldump


# Supported Notifications


- Email


# How does it work?


1. Download the [latest release](https://github.com/agorman/repbak/releases).
2. Create a YAML configuration file
3. Run it `repbak -conf repbak.yaml`


# Configuration file


The YAML file defines repbak's operation.


## Full config example

~~~
log_path: /var/log/repbak.log
log_level: error
lib_path: /var/lib/repbak
time_format: Mon Jan 02 03:04:05 PM MST
retention: 7
http:
  addr: 0.0.0.0
  port: 4060
mysqldump:
  retention: 30
  output_path: /mnt/backups/mysql.dump
  schedule: "0 0 * * *"
  executable_path: mysqldump
  executable_args: --add-drop-database --all-databases -u user -ppass -h 127.0.0.1
  time_limit: 8h
email:
  host: mail.me.com
  port: 587
  user: me
  pass: pass
  starttls: true
  ssl: false
  subject: Database Replication Failure
  from: me@me.com
  to:
    - you@me.com
  history_subject: Database Backup History
  history_schedule: "0 0 * * *"
  history_template: /home/repbak/email.template
  on_failure: true
~~~


## Global Options


**log_path** - File on disk where repbak logs will be stored. Defaults to /var/log/repbak.log.

**log_level** - Sets the log level. Valid levels are: panic, fatal, trace, debug, warn, info, and error. Defaults to error.

**lib_path** - The directory on disk where repbak lib files are stored. Defaults to /var/lib/repbak.

**time_format** - The format used when displaying backup stats. See formatting options in the go time.Time package. Defaults to Mon Jan 02 03:04:05 PM MST.

**retention** - The number of stats that are stored for each backup. If set to less than 0 no stats are saved. Defaults to 7.

## HTTP


**addr** - The listening address for the HTTP server. Default to 127.0.0.1

**port** - The listening port for the HTTP server. Default to 4040


## mysqldump


**retention** - The number of backups to keep before rotating old backups out. Defaults to 7.

**output_path** - The path where backups will be stored.

**schedule** - The cron expression that defines when backups are created.
    
**executable_path** - The path to the mysqldump binary. Defaults to mysqldump.

**executable_args** - The arguments passed to the executable used to create the mysql backup. Defaults to --add-drop-database --all-databases.
    
**time_limit** - Optional limit to the time it takes to run the backup.


## Email


**host** - The hostname or IP of the SMTP server.

**port** - The port of the SMTP server.

**user** - The username used to authenticate.

**pass** - The password used to authenticate.

**start_tls** - StartTLS enables TLS security. If both StartTLS and SSL are true then StartTLS will be used.

**insecure_skip_verify** - When using TLS skip verifying the server's certificate chain and host name.

**ssl** - SSL enables SSL security. If both StartTLS and SSL are true then StartTLS will be used.

**from** - The email address the email will be sent from.

**to** - An array of email addresses for which emails will be sent.

**history_subject** - An optional subject to use when sending sync history emails. Defaults to Database Backup History.

**history_schedule** - An optional cron expression. If set then an email with sync history will be sent based on the schedule.

**history_template** - 	An optional path to an email template to use when sending history emails. If not set uses the default template.

**on_failure** - An optional value that will send an email for each backup failure if true.


# Flags


**-conf** - Path to the repbak configuration file

**-debug** - Log to STDOUT


# HTTP Health Checks


The optional HTTP server creates two endpoints.

**/live** - A liveness check that always returns 200. 

**/health** - A health check that returns 200 if the latest run for each backup was successful and 503 otherwise.


## Road Map


- Docker Image
- Systemd service file
- Create rpm
- Create deb
- Support for more dumpers
- Support for more notifiers
