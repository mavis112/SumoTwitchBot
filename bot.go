package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

type bot struct {
	client     *twitch.Client
	httpClient *http.Client

	channel       string
	translitMap   map[rune]string
	rikishiList   []*rikishi
	commands      map[string]commandHandler
	mtx           sync.RWMutex
	currBasho     *bashoTime
	isModsOnly    bool
	userCooldowns map[string]time.Time
}

func newBot() *bot {
	channel := os.Getenv("CHANNEL")
	b := &bot{
		client: newClient(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		channel:       channel,
		translitMap:   translitMap,
		mtx:           sync.RWMutex{},
		isModsOnly:    false,
		userCooldowns: make(map[string]time.Time),
	}
	if err := b.newRikishis(); err != nil {
		log.Fatal(err)
	}
	b.registerCommands()
	return b
}

func newClient() *twitch.Client {
	botName := os.Getenv("BOT_NAME")
	rawToken := strings.TrimPrefix(os.Getenv("OAUTH_TOKEN"), "oauth:")
	oauthToken := "oauth:" + rawToken
	client := twitch.NewClient(botName, oauthToken)
	return client
}

func (b *bot) newRikishis() error {
	url := "https://www.sumo-api.com/api/rikishis?limit=1000"
	resp, err := b.httpClient.Get(url)
	if err != nil {
		log.Println("Could not fetch rikishis from API:", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("Bad status code")
		return errors.New("bad status code in http response")
	}

	limitBody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data rikishiUtil

	err = json.NewDecoder(limitBody).Decode(&data)
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
		b.handleMessage(msg)
	}
}

func (b *bot) handleMessage(msg twitch.PrivateMessage) {
	if len(msg.Message) == 0 || msg.Message[0] != '!' {
		return
	}

	b.mtx.RLock()
	modsOnly := b.isModsOnly
	b.mtx.RUnlock()

	if modsOnly && (!msg.User.IsBroadcaster && !msg.User.IsMod) {
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
	if !msg.User.IsBroadcaster && !msg.User.IsMod {

		b.mtx.RLock()
		lastTime := b.userCooldowns[msg.User.ID]
		b.mtx.RUnlock()
		if time.Since(lastTime) < 3*time.Second {
			return
		}
		b.mtx.Lock()
		b.userCooldowns[msg.User.ID] = time.Now()
		b.mtx.Unlock()
	}
	comm.handler(msg, comm.resp)
}
