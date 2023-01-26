package cli

import (
	"errors"
	"fmt"
	"io"
	"syscall"

	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/routeros.v2"
	"routeros-exporter/pkg/conf"
)

func newClient(raw *routeros.Client, target *conf.Target, close func()) *Client {
	return &Client{
		Client: raw,
		target: target,
		close:  close,
	}
}

type Client struct {
	*routeros.Client
	target *conf.Target
	close  func()
}

func (c *Client) Run(sentence ...string) (*routeros.Reply, error) {
	reply, err := c.Client.Run(sentence...)
	if err != nil && (errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE)) {
		c.close()
		return nil, fmt.Errorf("conn error: %w", err)
	}
	return reply, err
}

func (c *Client) Labels() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("name", c.target.Name),
		attribute.String("host", c.target.Host),
	}
}
