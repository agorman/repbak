package repbak

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// MySQLDumpDumper dumps a mysql backup to a file.
type MySQLDumpDumper struct {
	config  *Config
	running bool
	mu      sync.Mutex
	cancel  context.CancelFunc
}

// NewMySQLDumpDumper creates a NewMySQLDumpDumper.
func NewMySQLDumpDumper(config *Config) *MySQLDumpDumper {
	return &MySQLDumpDumper{
		config: config,
		mu:     sync.Mutex{},
	}
}

// Dump dumps the mysql data to a file based on the settings in config.
func (d *MySQLDumpDumper) Dump() Stat {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if d.config.MySQLDump.timeLimit != 0 {
		ctx, cancel = context.WithTimeout(ctx, d.config.MySQLDump.timeLimit)
		defer cancel()
	}

	stat := NewStat("mysqldump", d.config.TimeFormat)

	log.Info("Running: mysqldump")

	// check if already running
	d.mu.Lock()
	if d.running {
		stat.Skip = true
		log.Warn("MySQL Dumper: skipping because the previous scheduled dump is still running")
		d.mu.Unlock()
		return stat
	}

	d.running = true
	d.cancel = cancel
	d.mu.Unlock()

	stat = func(stat Stat) Stat {
		// Rotate the dump files
		logger := &lumberjack.Logger{
			Filename:   d.config.MySQLDump.OutputPath,
			MaxSize:    0,
			MaxBackups: d.config.MySQLDump.Retention,
		}
		if err := logger.Rotate(); err != nil {
			stat = stat.Finish(err)
			return stat
		}

		args := strings.Fields(d.config.MySQLDump.ExecutableArgs)

		cmd := exec.CommandContext(ctx, d.config.MySQLDump.ExecutablePath, args...)

		if err := os.MkdirAll(filepath.Dir(d.config.MySQLDump.OutputPath), 0644); err != nil {
			stat.Finish(fmt.Errorf("MySQL Dumper: failed to create dump directory %s: %v", filepath.Dir(d.config.MySQLDump.OutputPath), err))
			return stat
		}

		// write output into the dump file
		dump, err := os.Create(d.config.MySQLDump.OutputPath)
		if err != nil {
			stat.Finish(fmt.Errorf("MySQL Dumper: failed to create dump file %s: %v", d.config.MySQLDump.OutputPath, err))
			return stat
		}
		cmd.Stdout = dump

		stderr, err := cmd.StderrPipe()
		if err != nil {
			stat.Finish(fmt.Errorf("MySQL Dumper: failed to get STDERR pipe: %v", err))
			return stat

		}

		// write any errors into the log file
		scanner := bufio.NewScanner(stderr)
		go func() {
			// Read line by line and process it
			for scanner.Scan() {
				line := scanner.Text()
				log.Error(line)
			}
		}()

		// start the command
		if err := cmd.Start(); err != nil {
			stat.Finish(fmt.Errorf("MySQL Dumper: failed to start backup: %v", err))
			return stat
		}

		err = cmd.Wait()
		stat = stat.Finish(err)
		return stat
	}(stat)

	if stat.Success {
		log.Infof("Finished %s after %s", stat.Name, stat.Duration)
	} else {
		log.Errorf("Error %s: after %s: %s", stat.Name, stat.Duration, stat.Error)

	}

	d.mu.Lock()
	defer d.mu.Unlock()
	d.running = false
	d.cancel = nil

	return stat
}

// Stop stops the current dump if one is running.
func (d *MySQLDumpDumper) Stop() {
	if d.cancel != nil {
		d.cancel()
	}
}
