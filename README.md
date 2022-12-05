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


# HTTP Healthcheck


Repbak can optionally listen for HTTP liveness probes at /healthcheck. It will return a 200 status code if live.


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
~~~


## Global Options


**log_path** - File on disk where repbak logs will be stored. Defaults to /var/log/repbak.log.

**log_level** - Sets the log level. Valid levels are: panic, fatal, trace, debug, warn, info, and error. Defaults to error.


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


# Flags


**-conf** - Path to the repbak configuration file

**-debug** - Log to STDOUT


## Road Map


- Docker Image
- Systemd service file
- Create rpm
- Create deb
- Support for more dumpers
- Support for more notifiers