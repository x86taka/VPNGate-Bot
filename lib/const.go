package lib

import "time"

type Host struct {
	HostName string `gorm:"primary_key"`
	IP       string `gorm:"size:255"`
	PORT     int
}

type DbConfig struct {
	// データベース
	Dialect string

	// ユーザー名
	DBUser string

	// パスワード
	DBPass string

	// プロトコル
	DBProtocol string

	// DB名
	DBName string
}

type ServerList struct {
	HostName                  string        //0
	IP                        string        //1
	Score                     int           //2
	Ping                      int           //3
	Speed                     int           //4
	Country                   string        //5
	CountryShort              string        //6
	NumVpnSessions            int           //7
	Uptime                    time.Duration //8
	TotalUser                 string        //9
	TotalTraffic              int           //10
	LogType                   string        //11
	Operator                  string        //12
	Message                   string        //13
	OpenVPN_ConfigData_Base64 string        //14
}
