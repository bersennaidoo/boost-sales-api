package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/conf"
)

type Config struct {
	conf.Version
	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3000"`
		DebugHost       string        `conf:"default:0.0.0.0:4000"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:10s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
		ShutdownTimeout time.Duration `conf:"default:20s,mask"`
	}
	/*Auth struct {
		KeysFolder string `conf:"default:zarf/keys/"`
		ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
	}
	DB struct {
		User         string `conf:"default:postgres"`
		Password     string `conf:"default:postgres,mask"`
		Host         string `conf:"default:localhost"`
		Name         string `conf:"default:postgres"`
		MaxIdleConns int    `conf:"default:0"`
		MaxOpenConns int    `conf:"default:0"`
		DisableTLS   bool   `conf:"default:true"`
	}
	Zipkin struct {
		ReporterURI string  `conf:"default:http://localhost:9411/api/v2/spans"`
		ServiceName string  `conf:"default:sales-api"`
		Probability float64 `conf:"default:0.05"`
	}*/
}

func New(build string) (Config, error) {

	cfg := Config{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyright information here",
		},
	}

	const prefix = "BOOST-SALES"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
