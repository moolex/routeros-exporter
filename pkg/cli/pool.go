package cli

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/alphadose/haxmap"
	"golang.org/x/sync/singleflight"
	"gopkg.in/routeros.v2"
	"routeros-exporter/pkg/conf"
)

func New(def *conf.Defaults) *Pool {
	return &Pool{
		def: def,
		cli: haxmap.New[string, *Client](),
	}
}

type Pool struct {
	def *conf.Defaults
	cli *haxmap.Map[string, *Client]
	sfg singleflight.Group
}

func (p *Pool) Get(t *conf.Target) (*Client, error) {
	key := p.def.Username + "@" + t.Name
	if cli, has := p.cli.Get(key); has {
		return cli, nil
	}

	cli, err, _ := p.sfg.Do(key, func() (any, error) { return p.connect(t) })
	if err != nil {
		return nil, err
	}

	p.cli.Set(key, cli.(*Client))
	return cli.(*Client), err
}

func (p *Pool) connect(t *conf.Target) (*Client, error) {
	var (
		nc  net.Conn
		err error
	)
	if p.def.TLS {
		nc, err = tls.DialWithDialer(
			&net.Dialer{Timeout: p.def.Timeout},
			"tcp", fmt.Sprintf("%s:%d", t.Host, p.def.Port),
			&tls.Config{InsecureSkipVerify: p.def.Insecure},
		)
	} else {
		nc, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", t.Host, p.def.Port), p.def.Timeout)
	}

	if err != nil {
		return nil, err
	}

	rc, err := routeros.NewClient(nc)
	if err != nil {
		return nil, err
	}

	if err := rc.Login(p.def.Username, p.def.Password); err != nil {
		return nil, err
	}

	return newClient(rc, t), nil
}
