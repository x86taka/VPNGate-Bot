package main

import (
	"fmt"
	"github.com/x86taka/VPNGate-Bot/lib"
	"golang.org/x/xerrors"
	"math/rand"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	vcsession  *discordgo.VoiceConnection
	gate       *lib.VPNGate
	lastUpdate = time.Now()
)

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

	fmt.Println("Listening...")

	gate = lib.NewVPNGate()
	gate.Run()
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
		if h, err := RandomCountry(cn); err == nil {
			sendMessage(s, c, h)
		} else {
			errorMessage(s, c, err)
		}
		break
	case strings.HasPrefix(m.Content, "topsession"):
		cn := strings.Replace(m.Content, "topsession", "", 1)
		if h, err := TopSession(cn); err == nil {
			sendMessage(s, c, h)
		} else {
			errorMessage(s, c, err)
		}
	case strings.HasPrefix(m.Content, "topspeed"):
		cn := strings.Replace(m.Content, "topspeed", "", 1)
		if h, err := Top(cn); err == nil {
			sendMessage(s, c, h)
		} else {
			errorMessage(s, c, err)
		}
		break
	case m.Content == "ALL":
		if h, err := ALL(); err == nil {
			sendMessage(s, c, h)
		} else {
			errorMessage(s, c, err)
		}
		break
	}
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, lists *[]lib.ServerList) {
	var res string
	for _, list := range *lists {
		if p, ok := list.OriginalPing(); ok {
			res += fmt.Sprintf("```%s\nOriginal Ping: %s```", list.ToString(), p)
		} else {
			res += fmt.Sprintf("```%s\nHostDown!!```", list.ToString())
		}
	}

	_, err := s.ChannelMessageSend(c.ID, res+"\n LastUpdate "+gate.GetUpdateTime().Format("2006/1/2 15:04:05"))

	//log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func errorMessage(s *discordgo.Session, c *discordgo.Channel, err error) {

	_, err2 := s.ChannelMessageSend(c.ID, err.Error()+time.Now().Format("2006/1/2 15:04:05"))
	log.Println(">>> " + err.Error() + time.Now().Format("2006/1/2 15:04:05"))
	if err2 != nil {
		log.Println("Error sending message: ", err2)
	}
}

func Random() (*[]lib.ServerList, error) {
	servers := gate.GetServers()
	if len(servers) == 0 {
		return nil, xerrors.New("ServerList 0 Len")
	}
	return &[]lib.ServerList{servers[0]}, nil
}

func RandomCountry(cn string) (*[]lib.ServerList, error) {
	cn = strings.ToUpper(cn)
	count := 0
	var server lib.ServerList
	for _, server = range gate.GetServers() {
		if count >= 50 {
			return nil, xerrors.New("Error Not Found")
		}

		if server.CountryShort == cn {
			break
		}
		if strings.HasPrefix(strings.ToUpper(server.Country), cn) {
			break
		}
		count++
	}
	return &[]lib.ServerList{server}, nil
}

func Top(cn string) (*[]lib.ServerList, error) {
	servers := gate.GetServers()
	cn = strings.ToUpper(cn)
	top := 0
	sliceIndex := 0
	count := 0
	i := 0
	for i = 2; i < len(servers)-2; i++ {
		server := servers[i]
		if server.CountryShort == cn {
			if server.Speed > top {
				sliceIndex = i
				top = server.Speed
			}
		}
		if strings.HasPrefix(strings.ToUpper(server.Country), cn) {
			if server.Speed > top {
				sliceIndex = i
				top = server.Speed
			}
		}
		count++
	}
	if sliceIndex == 0 {
		return nil, xerrors.New("Error Not Found")
	}
	return &[]lib.ServerList{servers[sliceIndex]}, nil

}

func TopSession(cn string) (*[]lib.ServerList, error) {
	servers := gate.GetServers()
	cn = strings.ToUpper(cn)
	top := 0
	sliceIndex := 0
	count := 0
	i := 0
	for i = 2; i < len(servers)-2; i++ {
		server := servers[i]
		if server.CountryShort == cn {
			if server.NumVpnSessions > top {
				sliceIndex = i
				top = server.NumVpnSessions
			}
		}
		if strings.HasPrefix(strings.ToUpper(server.Country), cn) {
			if server.NumVpnSessions > top {
				sliceIndex = i
				top = server.NumVpnSessions
			}
		}
		count++
	}
	if sliceIndex == 0 {
		return nil, xerrors.New("Error Not Found")
	}
	return &[]lib.ServerList{servers[sliceIndex]}, nil
}

func ALL() (*[]lib.ServerList, error) {
	if len(gate.GetServers()) <= 10 || len(gate.GetServers()) <= 0 {
		return nil, xerrors.New("Error Not Found")
	}
	servers := gate.GetServers()[:10]
	return &servers, nil
}
