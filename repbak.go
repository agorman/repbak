package repbak

import (
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// RepBak reforms scheduled database backups.
type RepBak struct {
	config   *Config
	dumper   Dumper
	notifier Notifier
	crontab  *cron.Cron
	running  bool
	stopc    chan struct{}
	donec    chan struct{}
}

// New returns a new RepBak instance.
func New(config *Config, dumper Dumper, notifier Notifier) *RepBak {
	return &RepBak{
		config:   config,
		dumper:   dumper,
		notifier: notifier,
		crontab:  cron.New(),
		stopc:    make(chan struct{}),
		donec:    make(chan struct{}),
	}
}

// Start runs until stopped and does scheduled database backups.
func (r *RepBak) Start() error {
	if r.running {
		return nil
	}

	r.running = true

	log.Infof("Adding Schedule MySQL Backup: %s", r.config.MySQLDump.Schedule)
	// When more dumpers are added this can be generalized
	_, err := r.crontab.AddFunc(r.config.MySQLDump.Schedule, func() {
		log.Info("Dumping MySQL database")

		if err := r.backup(); err != nil {
			log.Errorf("Backup failed: %v", err)
		}
	})
	if err != nil {
		return err
	}

	go r.loop()

	return nil
}

// Stop stops repmon from running for schedule database backups.
func (r *RepBak) Stop() {
	if !r.running {
		return
	}

	r.stopc <- struct{}{}
	<-r.donec
	r.running = false
}

func (r *RepBak) loop() {
	r.crontab.Start()

	log.Infof("RepBak started")

	<-r.stopc

	r.crontab.Stop()
	r.dumper.Stop()
	r.donec <- struct{}{}
	log.Info("RepBak shutdown")
}

func (r *RepBak) backup() error {
	// Rotate the dump files
	logger := &lumberjack.Logger{
		Filename:   r.config.MySQLDump.OutputPath,
		MaxSize:    0,
		MaxBackups: r.config.MySQLDump.Retention,
	}
	if err := logger.Rotate(); err != nil {
		if err := r.notifier.Notify(err); err != nil {
			log.Errorf("Failed to send notification: %v", err)
		}
		return err
	}

	// Create a new dump
	if err := r.dumper.Dump(); err != nil {
		if err := r.notifier.Notify(err); err != nil {
			log.Errorf("Failed to send notification: %v", err)
		}
		return err
	}

	return nil
}
