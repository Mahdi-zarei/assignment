package config

import (
	"assignment/pkg/common"
	"bytes"
	"github.com/spf13/viper"
)

type PostgresConf struct {
	Host   string
	Port   int
	DBName string
}

type ListenConf struct {
	Host string
	Port string
}

type Config struct {
	GiftShopPGXConf PostgresConf
	EndpointConf    ListenConf
}

var defaultConf = `
GiftShopPGXConf:
	Host: 5.34.202.174
	Port: 5433
	DBName: giftshop
EndpointConf:
	Host: 0.0.0.0
	Port: 8080
`

func GetConfig() *Config {
	common.Must1(viper.ReadConfig(bytes.NewReader([]byte(defaultConf))))
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/config")
	viper.AddConfigPath("")

	var conf Config
	common.Must1(viper.Unmarshal(&conf))
	return &conf
}
