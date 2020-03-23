package lib

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getData() (string, error) {
	url := "https://www.vpngate.net/api/iphone/"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:58.0) Gecko/20100101 Firefox/58.0")

	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func GetIpPort(b64 string) []string {
	data, _ := base64.StdEncoding.DecodeString(b64)
	datas := string(data)
	datas = datas[strings.LastIndex(datas, "remote"):]
	ipport := strings.Split(datas[:strings.Index(datas, "\n")-1], " ")
	var re []string
	re = append(re, ipport[1])
	re = append(re, ipport[2])
	return re
}

func GetServers() []ServerList {
	res, err := getData()
	if err != nil {
		log.Println(err)
		return []ServerList{}
	}
	var result []ServerList
	csv := strings.Split(res, "\n")
	csv = csv[2 : len(csv)-2]
	for _, line := range csv {
		spline := strings.Split(line, ",")
		score, _ := strconv.Atoi(spline[2])
		ping, _ := strconv.Atoi(spline[3])
		speed, _ := strconv.Atoi(spline[4])
		NumVpnSessions, _ := strconv.Atoi(spline[7])
		TotalTraffic, _ := strconv.Atoi(spline[10])
		Uptime, _ := strconv.Atoi(spline[8])
		s := ServerList{
			HostName:                  spline[0],
			IP:                        spline[1],
			Score:                     score,
			Ping:                      ping,
			Speed:                     speed / 1000000,
			Country:                   spline[5],
			CountryShort:              spline[6],
			NumVpnSessions:            NumVpnSessions,
			Uptime:                    time.Duration(Uptime) * time.Millisecond,
			TotalUser:                 spline[9],
			TotalTraffic:              TotalTraffic / 1024 / 1024 / 1024 / 1024,
			LogType:                   spline[11],
			Operator:                  spline[12],
			Message:                   spline[13],
			OpenVPN_ConfigData_Base64: spline[14],
		}
		result = append(result, s)
	}
	return result
}

func (this *ServerList) ToString() string {
	return fmt.Sprintf(
		"HostName: %s.opengw.net\n"+
			"IP: %s\n"+
			"Score: %d\n"+
			"Ping: %d\n"+
			"Speed: %d\n"+
			"Country: %s\n"+
			"Sessions: %d\n"+
			"Uptime: %f Hour\n"+
			"TotalUser: %s\n"+
			"TotalTraffic: %dTB",
		this.HostName,
		this.IP,
		this.Score,
		this.Ping,
		this.Speed,
		this.CountryShort,
		this.NumVpnSessions,
		this.Uptime.Hours(),
		this.TotalUser,
		this.TotalTraffic,
	)
}
