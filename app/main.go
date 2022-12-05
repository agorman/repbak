package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/agorman/repbak"
	"github.com/etherlabsio/healthcheck/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	conf := flag.String("conf", "/etc/repbak.yaml", "Path to the repbak configuration file")
	debug := flag.Bool("debug", false, "Log to STDOUT")
	flag.Parse()

	config, err := repbak.OpenConfig(*conf)
	if err != nil {
		log.Fatal(err)
	}

	dumper := repbak.NewMySQLDumpDumper(config)

	notifier := repbak.NewEmailNotifier(config)

	if !*debug {
		logfile := &lumberjack.Logger{
			Filename:   config.LogPath,
			MaxSize:    1,
			MaxBackups: 10,
			MaxAge:     30,
		}
		log.SetOutput(logfile)
	}

	rb := repbak.New(config, dumper, notifier)
	rb.Start()
	defer rb.Stop()

	errc := make(chan error, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// TODO: add a healthcheck that makes sure tha last backup was successful
	if config.HTTP != nil {
		http.Handle("/healthcheck", healthcheck.Handler(
			healthcheck.WithChecker(
				"live", healthcheck.CheckerFunc(
					func(ctx context.Context) error {
						return nil
					},
				),
			),
		))

		go func() { errc <- http.ListenAndServe(fmt.Sprintf("%s:%d", config.HTTP.Addr, config.HTTP.Port), nil) }()
	}

	select {
	case s := <-sig:
		log.Warnf("Received signal %s, exiting", s)
		return
	case e := <-errc:
		log.Errorf("Run error: %s", e)
		return
	}
}
