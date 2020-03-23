package main

import (
	"fmt"
	"github.com/x86taka/VPNGate-Bot/lib"
	"math/rand"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"net"
	"os"
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

var csvdata []lib.ServerList

func Ping(ip string, port string) string {
	/* 	out, err := exec.Command("php", "-d", "display_errors=off", "script/ping.php", ip, port).Output()
	   	if err != nil {
	   		log.Println(err)
	   	}
		   return string(out) */
	str := ""
	current := 0
	var conn net.Conn = nil
	d := net.Dialer{Timeout: 2 * time.Second}
	for {
		var err error
		start := time.Now()
		if conn == nil {
			conn, err = d.Dial("tcp", ip+":"+port)
		}
		end := time.Now()
		if err != nil {
			str += "Host Down\n"
			fmt.Println(err)
		} else {
			str += strconv.FormatFloat((end.Sub(start)).Seconds()*1000.0, 'f', 2, 64)
			str += " ms"
		}
		conn = nil
		current++
		if current == 1 {
			break
		}
	}
	return str
}

func main() {
	rand.Seed(time.Now().UnixNano())
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
	csvdata = lib.GetServers()
	lastUpdate = time.Now()

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
	//fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, "random"):
		cn := strings.Replace(m.Content, "random", "", 1)
		sendMessage(s, c, RandomCountry(csvdata, cn))
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
		break
	}
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg+"\n LastUpdate "+lastUpdate.Format("2006/1/2 15:04:05"))

	//log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func Random(servers []lib.ServerList) string {
	var res string = ""
	if len(servers) == 0 {
		return "Error"
	}
	ipport := lib.GetIpPort(servers[0].OpenVPN_ConfigData_Base64)
	res = fmt.Sprintf("```%s\nOriginal Ping: %s```", servers[0].ToString(), Ping(ipport[0], ipport[1]))
	return res
}

func RandomCountry(servers []lib.ServerList, cn string) string {
	cn = strings.ToUpper(cn)
	var res string = ""
	count := 0
	var server lib.ServerList
	for _, server = range servers {
		if count >= 50 {
			return "Not Found"
		}

		if server.CountryShort == cn {
			break
		}
		if strings.HasPrefix(strings.ToUpper(server.Country), cn) {
			break
		}
		count++
	}
	ipport := lib.GetIpPort(server.OpenVPN_ConfigData_Base64)
	res = fmt.Sprintf("```%s\nOriginal Ping: %s```", server.ToString(), Ping(ipport[0], ipport[1]))
	return res
}

func Top(servers []lib.ServerList, cn string) string {
	cn = strings.ToUpper(cn)
	var res string = ""
	top := 0
	topindex := 0
	count := 0
	i := 0
	for i = 2; i < len(servers)-2; i++ {
		server := servers[i]
		if server.CountryShort == cn {
			if server.Speed > top {
				topindex = i
				top = server.Speed
			}
		}
		if strings.HasPrefix(strings.ToUpper(server.Country), cn) {
			if server.Speed > top {
				topindex = i
				top = server.Speed
			}
		}
		count++
	}
	if topindex == 0 {
		return "Not Found"
	}
	server := servers[topindex]
	ipport := lib.GetIpPort(server.OpenVPN_ConfigData_Base64)
	res = fmt.Sprintf("```%s\nOriginal Ping: %s```", server.ToString(), Ping(ipport[0], ipport[1]))
	return res
}

func TopSession(servers []lib.ServerList, cn string) string {
	cn = strings.ToUpper(cn)
	var res string = ""
	top := 0
	topindex := 0
	count := 0
	i := 0
	for i = 2; i < len(servers)-2; i++ {
		server := servers[i]
		if server.CountryShort == cn {
			if server.NumVpnSessions > top {
				topindex = i
				top = server.NumVpnSessions
			}
		}
		if strings.HasPrefix(strings.ToUpper(server.Country), cn) {
			if server.NumVpnSessions > top {
				topindex = i
				top = server.NumVpnSessions
			}
		}
		count++
	}
	if topindex == 0 {
		return "Not Found"
	}
	server := servers[topindex]
	ipport := lib.GetIpPort(server.OpenVPN_ConfigData_Base64)
	res = fmt.Sprintf("```%s\nOriginal Ping: %s```", server.ToString(), Ping(ipport[0], ipport[1]))
	return res
}

func ALL(servers []lib.ServerList) string {
	var res string = ""
	i := 0
	count := 0
	for i = 2; i < len(servers)-2; i++ {
		if count >= 10 {
			break
		}
		server := servers[i]
		ipport := lib.GetIpPort(server.OpenVPN_ConfigData_Base64)
		res += fmt.Sprintf("```%s\nOriginal Ping: %s```", server.ToString(), Ping(ipport[0], ipport[1]))
		count++
	}
	if res == "" {
		res = "Error...."
	}
	return res
}
