package conf

import (
	"time"
)

type Config struct {
	Global    *Global     `yaml:"global"`
	Defaults  *Defaults   `yaml:"defaults"`
	Targets   []*Target   `yaml:"targets"`
	Endpoints []*Endpoint `yaml:"endpoints"`
}

type Global struct {
	Listen   string        `yaml:"listen"`
	Interval time.Duration `yaml:"interval"`
}

type Defaults struct {
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Port     int           `yaml:"port"`
	TLS      bool          `yaml:"tls"`
	Insecure bool          `yaml:"insecure"`
	Timeout  time.Duration `yaml:"timeout"`
}

type Target struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type Endpoint struct {
	Path    string    `yaml:"path"`
	Metrics []*Metric `yaml:"metrics"`
	Labels  []*Label  `yaml:"labels"`
}

type Metric struct {
	Name string `yaml:"name"`
}

type Label struct {
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}
