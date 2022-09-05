package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HttpHost       string `mapstructure:"HTTP_HOST"`
	HttpPort       string `mapstructure:"HTTP_PORT"`
	MongoDBHost    string `mapstructure:"MONGODB_HOST"`
	MongoDBPort    string `mapstructure:"MONGODB_PORT"`
	MongoDBName    string `mapstructure:"MONGODB_NAME"`
	RpcDefaultHost string `mapstructure:"RPC_DEFAULT_HOST"`
	RpcPortUserApp string `mapstructure:"RPC_PORT_USER_APP"`
	RpcAuthService string `mapstructure:"RPC_AUTH_SERVICE"`
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
