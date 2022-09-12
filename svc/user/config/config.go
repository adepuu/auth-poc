package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HttpHost       string `mapstructure:"HTTP_HOST"`
	HttpPort       string `mapstructure:"HTTP_PORT"`
	MongoDBURI     string `mapstructure:"MONGODB_URI"`
	RpcAuthHost    string `mapstructure:"RPC_AUTH_HOST"`
	RpcAuthService string `mapstructure:"RPC_AUTH_SERVICE"`
	RpcDefaultHost string `mapstructure:"RPC_DEFAULT_HOST"`
	RpcPortUserApp string `mapstructure:"RPC_PORT_USER_APP"`
}

func LoadConfig(path, serverEnv string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(serverEnv)
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func (cfg Config) GetValues() Config {
	return cfg
}
