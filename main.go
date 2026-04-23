package main

import (
	"log"
	"os"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error: config.env file not found!")
	}

	bot := newBot()

	msgChan := make(chan twitch.PrivateMessage, 10)

	go bot.handleLoop(msgChan)

	bot.client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		msgChan <- msg
	})

	bot.client.OnConnect(func() {
		log.Printf("Successfully authorized as %s\n", os.Getenv("BOT_NAME"))
	})

	bot.client.OnSelfJoinMessage(func(msg twitch.UserJoinMessage) {
		log.Printf("Successfully joined channel: %s\n", msg.Channel)
		bot.client.Say(msg.Channel, "SumoTwitchBot connected! 🏯")
	})

	bot.client.Join(bot.channel)
	if err := bot.client.Connect(); err != nil {
		log.Fatal("Failed to join channel:", err)
	}
}
