package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type Config struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	Secret   string        `yaml:"secret" env-required:"true"`
	App      AppConfig
	Postgres PostgresConfig
	MC       MicrocontrollerConfig
	Log      LogConfig
}

type AppConfig struct {
	Port        uint16 `yaml:"port"`
	Addr        string `yaml:"addr"`
	MachineAddr string `yaml:"machine_local_addr"`
}

type PostgresConfig struct {
	Addr           string        `yaml:"addr"`
	Port           uint16        `yaml:"port"`
	User           string        `yaml:"user"`
	Password       string        `yaml:"password"`
	DB             string        `yaml:"db"`
	ConnTimeExceed time.Duration `yaml:"conn_time_exceed"`
}

type MicrocontrollerConfig struct {
	RequestTimeout time.Duration `yaml:"request_timeout" env-required:"true"`
}

type LogConfig struct {
	OutDir string `yaml:"out_dir"`
	Dev    string `yaml:"dev"`
	CSV    string `yaml:"csv"`
}

// Function will panic if can not read config file or environment variables
func MustLoad() *Config {
	var cfg Config

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(errors.Wrap(err, "config file dost not exists"))
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(errors.Wrap(err, "failed to read config"))
	}
	return &cfg
}

// fetchConfigPath fetches config path from command line flar or env variable
// Priority: flag > env > default
// Default value is empty string
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path fo config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
