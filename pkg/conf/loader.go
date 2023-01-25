package conf

import (
	"github.com/spf13/viper"
)

func New(file string) *Loader {
	v := viper.New()
	v.SetConfigFile(file)
	v.SetConfigType("yaml")
	return &Loader{v: v}
}

type Loader struct {
	v *viper.Viper
}

func (l *Loader) Reload() (*Config, error) {
	if err := l.v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := l.v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
