package config

import (
	"assignment/pkg/common"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
)

type PostgresConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func (p *PostgresConf) GenerateConnectURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", p.User, p.Password, p.Host, p.Port, p.DBName)
}

type ListenConf struct {
	Host string
	Port int
}

type Config struct {
	GiftShopPGXConf PostgresConf
	EndpointConf    ListenConf
}

var defaultConf = `
GiftShopPGXConf:
  Host: "5.34.202.174"
  Port: 5433
  User: "postgres"
  Password: "dummypass"
  DBName: "giftshop"
EndpointConf:
  Host: "0.0.0.0"
  Port: 8080
`

func GetConfig() *Config {
	viper.SetConfigType("yaml")
	common.Must1(viper.ReadConfig(bytes.NewReader([]byte(defaultConf))))
	viper.SetConfigName("conf")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("/")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("no config file found, using default")
	}

	var conf Config
	common.Must1(viper.Unmarshal(&conf))
	return &conf
}
