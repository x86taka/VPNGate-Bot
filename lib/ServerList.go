package lib

import (
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

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

func NewServerList() *ServerList {
	return &ServerList{}
}

func (l *ServerList) GetIpPort() (string, string) {
	data, _ := base64.StdEncoding.DecodeString(l.OpenVPN_ConfigData_Base64)
	datas := string(data)
	datas = datas[strings.LastIndex(datas, "remote"):]
	ipport := strings.Split(datas[:strings.Index(datas, "\n")-1], " ")
	return ipport[1], ipport[2]
}

func (l *ServerList) OriginalPing() (string, bool) {
	ip, port := l.GetIpPort()
	var str string
	var conn net.Conn = nil
	d := net.Dialer{Timeout: 1 * time.Second}
	var err error
	start := time.Now()
	if conn == nil {
		conn, err = d.Dial("tcp", ip+":"+port)
	}
	end := time.Now()
	if err != nil {
		return "", false
	} else {
		str = strconv.FormatFloat((end.Sub(start)).Seconds()*1000.0, 'f', 2, 64)
	}
	conn.Close()

	return str, true
}

func (l *ServerList) ToString() string {
	return fmt.Sprintf(
		"HostName: %s.opengw.net\n"+
			"IP: %s\n"+
			"Score: %d\n"+
			"Ping: %d\n"+
			"Speed: %d Mbps\n"+
			"Country: %s\n"+
			"Sessions: %d\n"+
			"Uptime: %f Hour\n"+
			"TotalUser: %s\n"+
			"TotalTraffic: %dTB",
		l.HostName,
		l.IP,
		l.Score,
		l.Ping,
		l.Speed,
		l.CountryShort,
		l.NumVpnSessions,
		l.Uptime.Hours(),
		l.TotalUser,
		l.TotalTraffic,
	)
}
