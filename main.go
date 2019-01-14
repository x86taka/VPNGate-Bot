package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	stopBot    = make(chan bool)
	vcsession  *discordgo.VoiceConnection
	lastUpdate = time.Now()
)

var csvdata []string

func getIpPort(b64 string) []string {
	data, _ := base64.StdEncoding.DecodeString(b64)
	datas := string(data)
	datas = datas[strings.LastIndex(datas, "remote"):]
	ipport := strings.Split(datas[:strings.Index(datas, "\n")], " ")
	var re []string
	re = append(re, ipport[1])
	re = append(re, ipport[2])
	return re
}

func Ping(ip string, port string) string {
	out, _ := exec.Command("php", "-d", "display_errors=off", "script/ping.php", ip, port).Output()
	return string(out)
}

func main() {

	content, err := ioutil.ReadFile("token.conf")
	if err != nil {
		log.Fatal(err)
	}

	discord, err := discordgo.New()
	discord.Token = string(content)
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
		os.Exit(1)
	}

	discord.AddHandler(onMessageCreate) //全てのWSAPIイベントが発生した時のイベントハンドラを追加
	// websocketを開いてlistening開始
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			csvdata = GetCSV()
			time.Sleep(10 * time.Minute)
			lastUpdate = time.Now()
		}
	}()

	fmt.Println("Listening...")
	<-stopBot //プログラムが終了しないようロック
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID) //チャンネル取得
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, "random"): //Bot宛に!helloworld コマンドが実行された時
		cn := strings.Replace(m.Content, "random", "", 1)
		sendMessage(s, c, RandomJP(csvdata, cn))
		break
	case strings.HasPrefix(m.Content, "topsession"):
		cn := strings.Replace(m.Content, "topsession", "", 1)
		sendMessage(s, c, TopSession(csvdata, cn))
		break
	case strings.HasPrefix(m.Content, "topspeed"):
		cn := strings.Replace(m.Content, "topspeed", "", 1)
		sendMessage(s, c, Top(csvdata, cn))
		break
	case m.Content == "ALL":
		sendMessage(s, c, ALL(csvdata))
		sendMessage(s, c, ALL(csvdata))
		break
	}
}

//メッセージを受信した時の、声の初めと終わりにPrintされるようだ
func onVoiceReceived(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	//log.Print("しゃべったあああああ")
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg+"\n LastUpdate "+lastUpdate.String())

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func GetCSV() []string {
	url := "https://www.vpngate.net/api/iphone/"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:58.0) Gecko/20100101 Firefox/58.0")

	//	dump, _ := httputil.DumpRequestOut(req, true)
	//	fmt.Printf("%s", dump)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	res := string(b)
	csv := strings.Split(res, "\n")
	defer resp.Body.Close()
	return csv[4 : len(csv)-3]
}

func Sort(s string, datas []string) string {
	var res string
	switch s {
	case "score":
		var i int
		for i = 2; i < 12; i++ {
			data := strings.Split(datas[i], ",")
			res += "Host : " + data[0]
			res += "\nSpeed : " + data[3]
			res += "\nCountry : " + data[6]
			res += "\nSessions : " + data[7] + "\n"
		}
		break
	case "speed":

	}

	return res
}

func Random(data []string) string {
	var res string = ""
	rand.Seed(time.Now().UnixNano())
	data1 := data[rand.Intn(30)+2]
	server := strings.Split(data1, ",")
	res += "Host : " + server[0]
	res += ".opengw.net\nIP : " + server[1]
	res += "\nSpeed : " + server[4]
	res += "\nCountry : " + server[6]
	res += "\nSessions : " + server[7] + "\n original Ping\n"
	ipport := getIpPort(server[14])
	res += Ping(ipport[0], ipport[1])
	return res
}

func RandomJP(data []string, cn string) string {
	cn = strings.ToUpper(cn)
	var res string = ""
	rand.Seed(time.Now().UnixNano())
	var server []string
	count := 0
	i := 0
	for i = 0; i < 3; i++ {
		for {
			if count >= 100 {
				return "Not Found"
			}
			data1 := data[rand.Intn(len(data)-2)+2]
			server = strings.Split(data1, ",")
			if server[6] == cn {
				break
			}
			if strings.HasPrefix(strings.ToUpper(server[5]), cn) {
				break
			}
			count++
		}
		res += "```Host : " + server[0]
		res += ".opengw.net\nIP : " + server[1]
		speed, _ := strconv.Atoi(server[4])
		res += "\nSpeed : " + strconv.Itoa(speed/1000000) + "Mbps"
		res += "\nCountry : " + server[5]
		res += "\nCountryShort : " + server[6]
		res += "\nSessions : " + server[7] + "\n original Ping\n"
		ipport := getIpPort(server[14])
		res += Ping(ipport[0], ipport[1]) + "```"
	}
	return res
}

func Top(data []string, cn string) string {
	cn = strings.ToUpper(cn)
	var res string = ""
	rand.Seed(time.Now().UnixNano())
	var server []string
	top := 0
	topindex := 0
	count := 0
	i := 0
	for i = 2; i < len(data)-2; i++ {
		data1 := data[i]
		server = strings.Split(data1, ",")
		if server[6] == cn {
			speed, _ := strconv.Atoi(server[4])
			if speed > top {
				topindex = i
				top = speed
			}
		}
		if strings.HasPrefix(strings.ToUpper(server[5]), cn) {
			speed, _ := strconv.Atoi(server[4])
			if speed > top {
				topindex = i
				top = speed
			}
		}
		count++
	}
	if topindex == 0 {
		return "Not Found"
	}
	server = strings.Split(data[topindex], ",")
	res += "```Host : " + server[0]
	res += ".opengw.net\nIP : " + server[1]
	speed, _ := strconv.Atoi(server[4])
	res += "\nSpeed : " + strconv.Itoa(speed/1000000) + "Mbps"
	res += "\nCountry : " + server[5]
	res += "\nCountryShort : " + server[6]
	res += "\nSessions : " + server[7] + "\n original Ping\n"
	ipport := getIpPort(server[14])
	res += Ping(ipport[0], ipport[1]) + "```"
	return res
}

func TopSession(data []string, cn string) string {
	cn = strings.ToUpper(cn)
	var res string = ""
	rand.Seed(time.Now().UnixNano())
	var server []string
	top := 0
	topindex := 0
	count := 0
	i := 0
	for i = 2; i < len(data)-2; i++ {
		data1 := data[i]
		server = strings.Split(data1, ",")
		if server[6] == cn {
			session, _ := strconv.Atoi(server[7])
			if session > top {
				topindex = i
				top = session
			}
		}
		if strings.HasPrefix(strings.ToUpper(server[5]), cn) {
			session, _ := strconv.Atoi(server[7])
			if session > top {
				topindex = i
				top = session
			}
		}
		count++
	}
	if topindex == 0 {
		return "Not Found"
	}
	server = strings.Split(data[topindex], ",")
	res += "```Host : " + server[0]
	res += ".opengw.net\nIP : " + server[1]
	speed, _ := strconv.Atoi(server[4])
	res += "\nSpeed : " + strconv.Itoa(speed/1000000) + "Mbps"
	res += "\nCountry : " + server[5]
	res += "\nCountryShort : " + server[6]
	res += "\nSessions : " + server[7] + "\n original Ping\n"
	ipport := getIpPort(server[14])
	res += Ping(ipport[0], ipport[1]) + "```"
	return res
}

func ALL(data []string) string {
	var res string = ""
	i := 0
	count := 0
	for i = 2; i < len(data)-2; i++ {
		if count >= 10 {
			break
		}
		data1 := data[i]
		server := strings.Split(data1, ",")
		res += "```Host : " + server[0]
		res += ".opengw.net\nIP : " + server[1]
		speed, _ := strconv.Atoi(server[4])
		res += "\nSpeed : " + strconv.Itoa(speed/1000000) + "Mbps"
		res += "\nCountry : " + server[5]
		res += "\nCountryShort : " + server[6]
		res += "\nSessions : " + server[7]
		ipport := getIpPort(server[14])
		res += "\nOpenVPN ip:Port : " + ipport[0] + ":" + ipport[1] + "```"
		count++
	}
	if res == "" {
		res = "Error...."
	}
	return res
}
