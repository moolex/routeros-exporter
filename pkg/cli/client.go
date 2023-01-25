package cli

import (
	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/routeros.v2"
	"routeros-exporter/pkg/conf"
)

func newClient(raw *routeros.Client, target *conf.Target) *Client {
	return &Client{Client: raw, target: target}
}

type Client struct {
	*routeros.Client
	target *conf.Target
}

func (c *Client) Labels() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("name", c.target.Name),
		attribute.String("host", c.target.Host),
	}
}
