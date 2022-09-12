package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HttpHost       string `mapstructure:"HTTP_HOST"`
	HttpPort       string `mapstructure:"HTTP_PORT"`
	MongoDBURI     string `mapstructure:"MONGODB_URI"`
	RpcDefaultHost string `mapstructure:"RPC_DEFAULT_HOST"`
	RpcPortAuthApp string `mapstructure:"RPC_PORT_AUTH_APP"`
	RpcUserHost    string `mapstructure:"RPC_USER_HOST"`
	RpcUserService string `mapstructure:"RPC_USER_SERVICE"`
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
