package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	"routeros-exporter/pkg/cli"
	"routeros-exporter/pkg/conf"
	"routeros-exporter/pkg/harvest"
	"routeros-exporter/pkg/log"
)

func main() {
	var config string
	flag.StringVar(&config, "c", "", "-c config.yaml")
	flag.Parse()

	logger, _ := zap.NewDevelopment()

	ld := conf.New(config)
	cfg, err := ld.Reload()
	if err != nil {
		logger.Fatal("load config failed", zap.Error(err))
	}

	pool := cli.New(cfg.Defaults)

	sch := gocron.NewScheduler(time.Local)
	if _, err := sch.Every(cfg.Global.Interval).Do(func() {
		harvest.Collect(
			log.WithLogger(context.TODO(), logger),
			pool,
			cfg.Targets,
			cfg.Endpoints,
		)
	}); err != nil {
		logger.Fatal("schedule failed", zap.Error(err))
	}

	sch.StartAsync()

	exporter, err := prometheus.New()
	if err != nil {
		logger.Fatal("create exporter failed", zap.Error(err))
	}
	global.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))

	http.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr: cfg.Global.Listen,
	}

	go func() { _ = server.ListenAndServe() }()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{}, 1)
	go func() {
		<-sigs
		_ = server.Shutdown(context.Background())
		done <- struct{}{}
	}()

	<-done
}
