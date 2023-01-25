package harvest

import (
	"context"
	"strings"

	"github.com/alphadose/haxmap"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
)

const (
	prefix = "routeros_"
)

var dm = newMeters()

func newMeters() *Meters {
	return &Meters{
		registry: haxmap.New[string, Gauge](),
	}
}

type Meters struct {
	registry *haxmap.Map[string, Gauge]
}

func (m *Meters) key(path, name string) string {
	return prefix + strings.ReplaceAll(path, "/", "_") + "_" + strings.ReplaceAll(name, "-", "_")
}

func (m *Meters) Gauge(path, name string) (Gauge, error) {
	key := m.key(path, name)
	if meter, has := m.registry.Get(key); has {
		return meter, nil
	}

	meter := global.Meter("routeros-exporter")

	gauge, err := meter.AsyncFloat64().Gauge(key)
	if err != nil {
		return nil, err
	}

	proxy := newGauge()

	if err := meter.RegisterCallback([]instrument.Asynchronous{gauge}, func(ctx context.Context) {
		for _, g := range proxy.gauges() {
			gauge.Observe(ctx, g.v, g.a...)
		}
	}); err != nil {
		return nil, err
	}

	m.registry.Set(key, proxy)
	return proxy, nil
}
