package lib

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type VPNGate struct {
	lastUpdate  time.Time
	serverLists *[]ServerList
	lock        sync.Mutex
}

func NewVPNGate() *VPNGate {
	return &VPNGate{}
}

func (v *VPNGate) Run() {
	t := time.NewTicker(15 * time.Minute)
	for {
		v.update()
		<-t.C
	}
}

func (v *VPNGate) GetServers() []ServerList {
	v.lock.Lock()
	v.lock.Unlock()
	return *v.serverLists
}

func (v *VPNGate) GetUpdateTime() time.Time {
	return v.lastUpdate
}

func (v *VPNGate) getHTTPRequest() (string, error) {
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

func (v *VPNGate) update() error {
	res, err := v.getHTTPRequest()
	v.lock.Lock()
	defer v.lock.Unlock()
	if err != nil {
		return err
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
	v.serverLists = &result
	v.lastUpdate = time.Now()
	return nil
}
