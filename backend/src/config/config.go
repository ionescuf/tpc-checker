package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Server Server
	GCP    GCP
}

type Server struct {
	Host string
}

type GCP struct {
	CredentialsFile string
	BucketName      string
	Token           string
}

func NewConfig() (*Config, error) {
	pflag.Parse()
	f := pflag.Lookup("config")
	configPath := ""
	if f == nil || f.Value.String() == "" {
		configPath = "./"
	} else {
		configPath = f.Value.String()
	}
	vp := viper.New()
	vp.AddConfigPath(configPath)
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.SetEnvPrefix("TPC_BE")
	vp.AutomaticEnv()
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = vp.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
