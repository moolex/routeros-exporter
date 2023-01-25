package harvest

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"routeros-exporter/pkg/cli"
	"routeros-exporter/pkg/conf"
	"routeros-exporter/pkg/log"
)

func Collect(ctx context.Context, pool *cli.Pool, targets []*conf.Target, endpoints []*conf.Endpoint) {
	logger := log.GetLogger(ctx)

	for _, target := range targets {
		client, err := pool.Get(target)
		if err != nil {
			logger.Warn("get client failed", zap.Error(err))
			continue
		}

		go func() {
			for _, endpoint := range endpoints {
				if err := collect(logger, client, endpoint); err != nil {
					logger.Warn("collect metrics failed", zap.Error(err))
				}
			}
		}()
	}
}

func collect(log *zap.Logger, cli *cli.Client, ep *conf.Endpoint) error {
	var (
		props   []string
		labels  = make(map[string]*conf.Label)
		metrics = make(map[string]*conf.Metric)
	)

	for _, label := range ep.Labels {
		props = append(props, label.Name)
		labels[label.Name] = label
	}

	for _, metric := range ep.Metrics {
		props = append(props, metric.Name)
		metrics[metric.Name] = metric
	}

	props = lo.Uniq(props)

	reply, err := cli.Run(fmt.Sprintf("/%s/print", ep.Path), "=.proplist="+strings.Join(props, ","))
	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	for _, ret := range reply.Re {
		attrs := cli.Labels()

		for key, label := range labels {
			val, has := ret.Map[key]
			if !has {
				continue
			}
			if label.Alias != "" {
				key = label.Alias
			}
			attrs = append(attrs, attribute.String(key, val))
		}

		for key := range metrics {
			val, has := ret.Map[key]
			if !has {
				continue
			}

			gauge, err := dm.Gauge(ep.Path, key)
			if err != nil {
				return fmt.Errorf("get meter failed: %w", err)
			}

			if f, err := strconv.ParseFloat(val, 64); err != nil {
				log.Info("parse as float failed", zap.Error(err))
			} else {
				gauge.Set(f, attrs...)
			}
		}
	}

	return nil
}
