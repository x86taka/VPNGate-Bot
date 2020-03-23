package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/x86taka/VPNGate-Bot/lib"
	"log"
	"strconv"
	"time"
)

func connectGorm(dbConfig lib.DbConfig) *gorm.DB {
	connectTemplate := "%s:%s@%s/%s?parseTime=true"
	connect := fmt.Sprintf(connectTemplate, dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBProtocol, dbConfig.DBName)
	db, err := gorm.Open(dbConfig.Dialect, connect)

	if err != nil {
		log.Println(err.Error())
	}

	return db
}

func main() {
	var db *gorm.DB
	if dbconfig, ok := lib.Config(); ok {
		db = connectGorm(dbconfig)
		db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&lib.Host{})
		db.LogMode(false)
	} else {
		log.Panic("Config NotFound")
	}

	var all []lib.Host
	db.Find(&all)
	for _, host := range all {
		fmt.Println(host.HostName + "	" + host.IP)
	}
	for {
		for _, server := range lib.GetServers() {

			IP_Port := lib.GetIpPort(server.OpenVPN_ConfigData_Base64)
			host := &lib.Host{}
			host.HostName = server.HostName
			if db.First(&host).RecordNotFound() {
				host.IP = IP_Port[0]
				port, _ := strconv.Atoi(IP_Port[1])
				host.PORT = port
				db.Create(&host)
			} else {
				host.IP = IP_Port[0]
				port, _ := strconv.Atoi(IP_Port[1])
				host.PORT = port
				db.Update(&host)
			}
		}
		time.Sleep(2 * time.Minute)
	}

}
