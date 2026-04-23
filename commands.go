package main

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

type commandHandler struct {
	handler func(msg twitch.PrivateMessage, resp response)
	resp    response
}

func (b *bot) registerCommands() {
	b.commands = map[string]commandHandler{
		"!stats":   {b.showRikishiStats, enResp},
		"!стат":    {b.showRikishiStats, ruResp},
		"!matchup": {b.showMatchUp, enResp},
		"!матчап":  {b.showMatchUp, ruResp},
		"!last":    {b.showLastMatches, enResp},
		"!ласт":    {b.showLastMatches, ruResp},
		"!next":    {b.showNextMatch, enResp},
		"!след":    {b.showNextMatch, ruResp},
		"!top5":    {b.showTop5Bouts, enResp},
		"!топ5":    {b.showTop5Bouts, ruResp},
	}
}

func (b *bot) showRikishiStats(msg twitch.PrivateMessage, resp response) {
	msgSlice := strings.Fields(msg.Message)
	if len(msgSlice) < 2 {
		b.client.Say(b.channel, fmt.Sprintf("@%s %s", msg.User.DisplayName, resp.ErrNoRikishi))
		return
	}
	rikishi := msgSlice[1]
	if resp.NeedTranslit {
		rikishi = translitShikona(rikishi, b.translitMap)
	}
	finalAnswer := getSingleStat(rikishi, msg.User.DisplayName, b.httpClient, b.rikishiList, resp)

	b.client.Say(b.channel, finalAnswer)
}

func (b *bot) showMatchUp(msg twitch.PrivateMessage, resp response) {
	msgSlice := strings.Fields(msg.Message)
	if len(msgSlice) < 3 {
		b.client.Say(b.channel, fmt.Sprintf("@%s %s", msg.User.DisplayName, resp.ErrNoRikishi))
		return
	}
	rikishi1 := msgSlice[1]
	rikishi2 := msgSlice[2]

	if resp.NeedTranslit {
		rikishi1 = translitShikona(rikishi1, b.translitMap)
		rikishi2 = translitShikona(rikishi2, b.translitMap)
	}
	finalAnswer := GetMatchup(rikishi1, rikishi2, msg.User.DisplayName, b.httpClient, b.rikishiList, resp)
	b.client.Say(b.channel, finalAnswer)
}

func (b *bot) showLastMatches(msg twitch.PrivateMessage, resp response) {
	msgSlice := strings.Fields(msg.Message)
	if len(msgSlice) < 2 {
		b.client.Say(b.channel, fmt.Sprintf("@%s %s", msg.User.DisplayName, resp.ErrNoRikishi))
		return
	}
	rikishi := msgSlice[1]
	if resp.NeedTranslit {
		rikishi = translitShikona(rikishi, b.translitMap)
	}
	finalAnswer := getLastMatches(rikishi, msg.User.DisplayName, b.httpClient, b.rikishiList, resp)
	b.client.Say(b.channel, finalAnswer)
}

func (b *bot) showNextMatch(msg twitch.PrivateMessage, resp response) {
	msgSlice := strings.Fields(msg.Message)
	if len(msgSlice) < 2 {
		b.client.Say(b.channel, fmt.Sprintf("@%s %s", msg.User.DisplayName, resp.ErrNoRikishi))
		return
	}
	rikishi := msgSlice[1]
	if resp.NeedTranslit {
		rikishi = translitShikona(rikishi, b.translitMap)
	}
	finalAnswer := b.getNextMatch(rikishi, msg.User.DisplayName, resp)
	b.client.Say(b.channel, finalAnswer)
}

func (b *bot) showTop5Bouts(msg twitch.PrivateMessage, resp response) {
	finalAnswer := b.getTop5Bouts(msg.User.DisplayName, resp)
	b.client.Say(b.channel, finalAnswer)
}
