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
func (d *MySQLDumpDumper) Dump() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if d.config.MySQLDump.timeLimit != 0 {
		ctx, cancel = context.WithTimeout(ctx, d.config.MySQLDump.timeLimit)
		defer cancel()
	}

	// check if already running
	d.mu.Lock()
	if d.running {
		log.Warn("MySQL Dumper: skipping because the previous scheduled dump is still running")
		d.mu.Unlock()
		return nil
	}

	d.running = true
	d.cancel = cancel
	d.mu.Unlock()

	args := strings.Fields(d.config.MySQLDump.ExecutableArgs)

	cmd := exec.CommandContext(ctx, d.config.MySQLDump.ExecutablePath, args...)

	if err := os.MkdirAll(filepath.Base(d.config.MySQLDump.OutputPath), 0644); err != nil {
		return fmt.Errorf("MySQL Dumper: failed to create dump directory %s: %v", filepath.Base(d.config.MySQLDump.OutputPath), err)
	}

	// write output into the dump file
	dump, err := os.Create(d.config.MySQLDump.OutputPath)
	if err != nil {
		return fmt.Errorf("MySQL Dumper: failed to create dump file %s: %v", d.config.MySQLDump.OutputPath, err)
	}
	cmd.Stdout = dump

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("MySQL Dumper: failed to get STDERR pipe: %v", err)

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
		return fmt.Errorf("MySQL Dumper: failed to start backup: %v", err)

	}

	err = cmd.Wait()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.running = false
	d.cancel = nil

	return err
}

// Stop stops the current dump if one is running.
func (d *MySQLDumpDumper) Stop() {
	if d.cancel != nil {
		d.cancel()
	}
}
