package repbak

import (
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// RepBak reforms scheduled database backups.
type RepBak struct {
	config   *Config
	db       DB
	dumper   Dumper
	notifier Notifier
	crontab  *cron.Cron
	running  bool
	stopc    chan struct{}
	donec    chan struct{}
}

// New returns a new RepBak instance.
func New(config *Config, db DB, dumper Dumper, notifier Notifier) *RepBak {
	return &RepBak{
		config:   config,
		db:       db,
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

	log.Infof("Adding Schedule For mysqldump: %s", r.config.MySQLDump.Schedule)
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

	// setup scheduled stats email
	if r.config.Email != nil && r.config.Email.HistorySchedule != "" {
		_, err := r.crontab.AddFunc(r.config.Email.HistorySchedule, func() {
			statMap, err := r.db.List()
			if err != nil {
				log.Error(err)
				return
			}

			if err := r.notifier.NotifyHistory(statMap); err != nil {
				log.Error(err)
			}
		})
		if err != nil {
			return err
		}

		log.Infof("History Email Scheduled: %s", r.config.Email.HistorySchedule)
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
	// Create a new dump
	stat := r.dumper.Dump()
	if stat.Skip {
		return nil
	}

	if !stat.Success && r.config.Email != nil && r.config.Email.OnFailure {
		if err := r.notifier.Notify(stat); err != nil {
			log.Error(err)
		}
	}

	if r.config.Retention > -1 {
		if err := r.db.Insert(stat); err != nil {
			log.Errorf("Failed to write stats for %s: %v", stat.Name, err)
		}
	}

	return stat.Error
}
