package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
"github.com/go-resty/resty/v2"
	"github.com/bwmarrin/discordgo"
	"time"
)

var Client *resty.Client

type Data struct {
	AssetIdBase string `json:"asset_id_base"`
	Rates       []struct {
		Time         time.Time `json:"time"`
		AssetIdQuote string    `json:"asset_id_quote"`
		Rate         float64   `json:"rate"`
	} `json:"rates"`
}

func main() {

	ConfigInit()

    Client = resty.New()

	dg, err := discordgo.New("Bot " + Conf.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(CreateOwnInvite)
	dg.AddHandler(BtcPrice)
	dg.AddHandler(EmbedExample)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func CreateOwnInvite(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == Conf.BotPrefix+ "invite" {
		for _, guild := range s.State.Guilds {

				fmt.Println(guild.MemberCount)

		}
		invite ,_ := s.ChannelInviteCreate(m.ChannelID, discordgo.Invite{
			MaxAge: 500,
			MaxUses: 10,
			Temporary: true,
			Unique: true,

		})

		s.ChannelMessageSend(m.ChannelID, "https://discord.gg/"+invite.Code)
	}
}
func BtcPrice(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == Conf.BotPrefix + "BTC" {
		resp, _ := Client.R().
			SetHeader("X-CoinAPI-Key", Conf.CryptoApi).
			Get("https://rest.coinapi.io/v1/exchangerate/BTC?invert=false")

		tempData := &Data{}
		json.Unmarshal(resp.Body(), tempData)
		for _, rate := range tempData.Rates {
			if rate.AssetIdQuote == "USD" {
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       "BTC Price",
					Description: fmt.Sprintf("$%f", rate.Rate),
					Image: &discordgo.MessageEmbedImage{
						URL: "https://www.cryptocompare.com/media/19633/btc.png",
					},
					Color:       0x00ff00,
				})
				return
			}
		}
	}

}


func EmbedExample(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "Hello" {
		_ , err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Hello",
			Description: "Hello, " + m.Author.Username + "!",
			Image: &discordgo.MessageEmbedImage{
				URL: "https://i.imgur.com/w3duR07.png",
			},
			Color: 	0x00ff00,
		})
		if err != nil {
			fmt.Println("error sending message,", err)
			return
		}
	}

}
