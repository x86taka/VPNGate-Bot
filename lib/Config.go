package lib

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func Config() (DbConfig, bool) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	f, err := os.Open("config.yaml")
	if err != nil {
		dbconfig := DbConfig{
			Dialect:    "mysql",
			DBUser:     "root",
			DBPass:     "",
			DBProtocol: "tcp(127.0.0.1:3306)",
			DBName:     "vpngate",
		}
		viper.SetDefault("DB", dbconfig)

		viper.WriteConfigAs("config.yaml")
	}
	if err := viper.ReadConfig(f); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			// Config file not found; ignore error if desired
			return DbConfig{}, false
		} else {
			// Config file was found but another error was produced
			log.Panicf("Config Error!!")
		}
	}
	m := viper.GetStringMapString("DB")

	return DbConfig{
		Dialect:    m["dialect"],
		DBUser:     m["dbuser"],
		DBPass:     m["dbpass"],
		DBProtocol: m["dbprotocol"],
		DBName:     m["dbname"],
	}, true
}
