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

var (
	gate *lib.VPNGate
	db   *gorm.DB
)

func connectGorm(dbConfig lib.DbConfig) error {
	connectTemplate := "%s:%s@%s/%s?parseTime=true"
	connect := fmt.Sprintf(connectTemplate, dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBProtocol, dbConfig.DBName)
	var err error
	db, err = gorm.Open(dbConfig.Dialect, connect)

	if err != nil {
		return err
	}
	return nil
}

func main() {
	var db *gorm.DB
	if dbconfig, ok := lib.Config(); ok {
		if err := connectGorm(dbconfig); err != nil {
			panic(err)
		}
		db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&lib.Host{})
		db.LogMode(false)
	} else {
		log.Panic("Config NotFound!! Please Edit config.yaml")
	}

	var all []lib.Host
	db.Find(&all)
	for _, host := range all {
		fmt.Println(host.HostName + "	" + host.IP)
	}
	gate = lib.NewVPNGate()
	go gate.Run()
	for {
		t := time.NewTicker(15 * time.Minute)
		for {
			insert()
			<-t.C
		}
	}
}

func insert() {
	for _, server := range gate.GetServers() {

		IP, Port := server.GetIpPort()
		host := &lib.Host{}
		host.HostName = server.HostName
		if db.First(&host).RecordNotFound() {
			host.IP = IP
			port, _ := strconv.Atoi(IP)
			host.PORT = port
			db.Create(&host)
		} else {
			host.IP = IP
			port, _ := strconv.Atoi(Port)
			host.PORT = port
			db.Update(&host)
		}
	}
	time.Sleep(2 * time.Minute)
}
