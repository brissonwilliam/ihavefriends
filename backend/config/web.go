package config

import "github.com/spf13/viper"

type Web struct {
	Address string
	JwtKey string
}

func GetWeb() Web {
	return Web{
		Address: viper.GetString("web_address"),
		JwtKey:  viper.GetString("web_jwtkey"),
	}
}