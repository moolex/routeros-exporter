package harvest

import (
	"github.com/alphadose/haxmap"
	"go.opentelemetry.io/otel/attribute"
)

type Gauge interface {
	Set(x float64, attrs ...attribute.KeyValue)
}

func newGauge() *float64Gauge {
	return &float64Gauge{
		g: haxmap.New[string, *f64Gauge](),
	}
}

type float64Gauge struct {
	g *haxmap.Map[string, *f64Gauge]
}

type f64Gauge struct {
	v float64
	a []attribute.KeyValue
}

func (f *float64Gauge) gauges() []*f64Gauge {
	var gs []*f64Gauge
	f.g.ForEach(func(_ string, value *f64Gauge) bool {
		gs = append(gs, value)
		return true
	})
	return gs
}

func (f *float64Gauge) Set(x float64, a ...attribute.KeyValue) {
	var k string
	for _, a2 := range a {
		k += string(a2.Key) + "/" + a2.Value.AsString() + ";"
	}
	v, _ := f.g.GetOrSet(k, &f64Gauge{v: x, a: a})
	v.v = x
}
