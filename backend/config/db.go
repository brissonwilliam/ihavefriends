package config

import (
	"github.com/spf13/viper"
)

type DB struct {
	DatabaseName    string
	Host            string
	Username        string
	Password        string
	Port            int
	TLS             bool
}

func GetDB() DB {
	return DB {
		DatabaseName:    viper.GetString("db_name"),
		Host:            viper.GetString("db_host"),
		Username:        viper.GetString("db_username"),
		Password:        viper.GetString("db_password"),
		Port:            viper.GetInt("db_port"),
		TLS:             viper.GetBool("db_tls"),
	}
}