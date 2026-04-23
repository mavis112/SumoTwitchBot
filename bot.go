package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

type bot struct {
	client     *twitch.Client
	httpClient *http.Client

	channel     string
	translitMap map[rune]string
	rikishiList []*rikishi
	commands    map[string]commandHandler
	isMonthEven bool
	currBasho   *bashoTime
}

func newBot() *bot {
	channel := os.Getenv("CHANNEL")
	b := &bot{
		client: newClient(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		channel:     channel,
		translitMap: translitMap,
		isMonthEven: isMonthEven(),
	}
	if err := b.newRikishis(); err != nil {
		log.Fatal(err)
	}
	b.registerCommands()
	return b
}

func newClient() *twitch.Client {
	botName := os.Getenv("BOT_NAME")
	oauthToken := "oauth:" + os.Getenv("OAUTH_TOKEN")
	client := twitch.NewClient(botName, oauthToken)
	return client
}

func (b *bot) newRikishis() error {
	url := "https://sumo-api.com/api/rikishis?limit=1000"
	resp, err := b.httpClient.Get(url)
	if err != nil {
		log.Println("Could not fetch rikishis from API:", err)
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Records []*rikishi `json:"records"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println("Could not decode rikishis API response:", err)
		return err
	}
	b.rikishiList = data.Records
	log.Printf("Total rikishis loaded: %d", len(b.rikishiList))
	return nil
}

func (b *bot) handleLoop(msgChan <-chan twitch.PrivateMessage) {
	for msg := range msgChan {
		go b.handleMessage(msg)
	}
}

func (b *bot) handleMessage(msg twitch.PrivateMessage) {
	if len(msg.Message) == 0 || msg.Message[0] != '!' {
		return
	}

	msgSlice := strings.Fields(msg.Message)
	if len(msgSlice) == 0 {
		return
	}

	command := strings.ToLower(msgSlice[0])

	comm, ok := b.commands[command]
	if !ok {
		return
	}

	comm.handler(msg, comm.resp)
}

func isMonthEven() bool {
	now := time.Now().UTC()

	return int(now.Month())%2 == 0
}
