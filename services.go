package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rikishiUtil struct {
	Records []*rikishi `json:"records"`
}
type rikishi struct {
	ID          int       `json:"id"`
	ShikonaEn   string    `json:"shikonaEn"`
	ShikonaJp   string    `json:"shikonaJp"`
	CurrentRank string    `json:"currentRank"`
	Heya        string    `json:"heya"`
	BirthDate   time.Time `json:"birthDate"`
	Shusshin    string    `json:"shusshin"`
	Height      float64   `json:"height"`
	Weight      float64   `json:"weight"`
	Debut       string    `json:"debut"`
}

type rikishiRank struct {
	RankValue int    `json:"rankValue"`
	Rank      string `json:"rank"`
}
type statsById struct {
	Basho           int `json:"basho"`
	TotalMatches    int `json:"totalMatches"`
	TotalWins       int `json:"totalWins"`
	Yusho           int `json:"yusho"`
	YushoByDivision struct {
		Jonidan   int `json:"Jonidan"`
		Jonokuchi int `json:"Jonokuchi"`
		Juryo     int `json:"Juryo"`
		Makushita int `json:"Makushita"`
		Makuuchi  int `json:"Makuuchi"`
		Sandanme  int `json:"Sandanme"`
	}
}

type matchUp struct {
	RikishiWins  int     `json:"rikishiWins"`
	OpponentWins int     `json:"opponentWins"`
	Matches      []match `json:"matches"`
}

type match struct {
	BashoID     string `json:"bashoId"`
	Division    string `json:"division"`
	Day         int    `json:"day"`
	MatchNo     int    `json:"matchNo"`
	EastID      int    `json:"eastId"`
	EastShikona string `json:"eastShikona"`
	EastRank    string `json:"eastRank"`
	WestID      int    `json:"westId"`
	WestShikona string `json:"westShikona"`
	WestRank    string `json:"westRank"`
	Kimarite    string `json:"kimarite"`
	WinnerID    int    `json:"winnerId"`
	WinnerEn    string `json:"winnerEn"`
	WinnerJp    string `json:"winnerJp"`
}

type matchUtil struct {
	Records []match `json:"records"`
}

type bashoTime struct {
	Date      string    `json:"date"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type torikumi struct {
	Torikumi []match `json:"torikumi"`
}

func getHighestRank(id int, httpClient *http.Client) string {
	url := fmt.Sprintf("https://www.sumo-api.com/api/ranks?rikishiId=%d", id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	limitBody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data []rikishiRank
	err = json.NewDecoder(limitBody).Decode(&data)
	if err != nil {
		return ""
	}
	minValue := 9999
	var rank string
	for _, rikishi := range data {
		if rikishi.RankValue < minValue {
			minValue = rikishi.RankValue
			rank = rikishi.Rank
		}
		if minValue == 1 {
			break
		}
	}
	return rank
}

func getStatsById(id int, httpClient *http.Client) (*statsById, error) {
	url := fmt.Sprintf("https://www.sumo-api.com/api/rikishi/%d/stats", id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not found")
	}
	limitBody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data statsById
	err = json.NewDecoder(limitBody).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func getMatchupById(id1, id2 int, httpClient *http.Client) (*matchUp, error) {
	url := fmt.Sprintf("https://www.sumo-api.com/api/rikishi/%d/matches/%d", id1, id2)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not found")
	}

	limitBody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data matchUp
	err = json.NewDecoder(limitBody).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (b *bot) getBashoDayById(r response) (int, string) {
	id := getBashoId()
	if id == "" {
		return 0, "even month"
	}
	b.mtx.RLock()
	if b.currBasho != nil && id == b.currBasho.Date {
		day, check := calculateDay(b.currBasho, r)
		b.mtx.RUnlock()
		return day, check
	}
	b.mtx.RUnlock()
	url := fmt.Sprintf("https://www.sumo-api.com/api/basho/%s", id)
	resp, err := b.httpClient.Get(url)
	if err != nil {
		return 0, "not ok"
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, "not ok"
	}
	limitBody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data bashoTime
	err = json.NewDecoder(limitBody).Decode(&data)
	if err != nil {
		return 0, "not ok"
	}
	b.mtx.Lock()
	if b.currBasho == nil || b.currBasho.Date != id {
		b.currBasho = &data
	}
	b.mtx.Unlock()
	day, check := calculateDay(&data, r)
	return day, check
}

func getTorikumiByDayId(id string, day int, httpClient *http.Client, r response) (torikumi, string) {
	url := fmt.Sprintf("https://www.sumo-api.com/api/basho/%s/torikumi/Makuuchi/%d", id, day)
	resp, err := httpClient.Get(url)
	if err != nil {
		return torikumi{}, r.ErrTechnical
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return torikumi{}, r.ErrTechnical
	}
	limitbody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data torikumi
	err = json.NewDecoder(limitbody).Decode(&data)
	if err != nil {
		return torikumi{}, r.ErrTechnical
	}
	if len(data.Torikumi) == 0 {
		return torikumi{}, r.ErrNoScheduleYet
	}
	return data, ""
}

func getSingleStat(name, user string, httpClient *http.Client, rikishisList []*rikishi, resp response) string {
	rikishi := getbyShikonaEn(name, rikishisList)
	if rikishi == nil {
		return fmt.Sprintf("@%s %s", user, resp.ErrNotFoundShikona)
	}
	ageStr := resp.ErrBadFormatDebut
	if !rikishi.BirthDate.IsZero() {
		ageStr = fmt.Sprintf("%d", int(time.Since(rikishi.BirthDate).Hours()/24/365.25))
	}

	var debut string
	if len(rikishi.Debut) == 6 {
		debut = rikishi.Debut[:4] + "." + rikishi.Debut[4:]
	} else {
		debut = resp.ErrBadFormatDebut
	}

	var (
		stats    *statsById
		err      error
		highRank string
	)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		stats, err = getStatsById(rikishi.ID, httpClient)
	}()
	go func() {
		defer wg.Done()
		highRank = getHighestRank(rikishi.ID, httpClient)
	}()
	wg.Wait()
	finalAnswer := fmt.Sprintf("@%s %s | %s: %s ", user, rikishi.ShikonaEn, resp.CurrentRank, rikishi.CurrentRank)
	if highRank != "" {
		finalAnswer = finalAnswer + fmt.Sprintf("| %s: %s ", resp.HighestRank, highRank)
	}
	finalAnswer = finalAnswer + fmt.Sprintf("| %s: %s | %d %s / %d %s | %s: %s | %s: %s ", resp.Heya, rikishi.Heya, int(rikishi.Height), resp.Height, int(rikishi.Weight), resp.Weight, resp.Age, ageStr, resp.Debut, debut)
	if err == nil {
		finalAnswer = finalAnswer + fmt.Sprintf("| %s: %d | %s: %d | %s: %d | %s: %d ", resp.Matches, stats.TotalMatches, resp.Wins, stats.TotalWins, resp.NumOfBasho, stats.Basho, resp.Yusho, stats.Yusho)
		if stats.Yusho != 0 {
			finalAnswer = finalAnswer + fmt.Sprintf("(%s: %d, %s: %d)", resp.Makuuchi, stats.YushoByDivision.Makuuchi, resp.Juryo, stats.YushoByDivision.Juryo)
		}
	}
	return finalAnswer
}

func GetMatchup(name1, name2, user string, httpClient *http.Client, rikishisList []*rikishi, resp response) string {
	errAnswer := fmt.Sprintf("@%s %s", user, resp.ErrNotFoundShikona)
	rikishi1 := getbyShikonaEn(name1, rikishisList)
	if rikishi1 == nil {
		return errAnswer
	}
	rikishi2 := getbyShikonaEn(name2, rikishisList)
	if rikishi2 == nil {
		return errAnswer
	}
	stat, err := getMatchupById(rikishi1.ID, rikishi2.ID, httpClient)
	if err != nil {
		return fmt.Sprintf("@%s %s", user, resp.ErrTechnical)
	}
	if stat.OpponentWins == 0 && stat.RikishiWins == 0 {
		return fmt.Sprintf("@%s %s", user, resp.ErrNoMatchup)
	}
	finalAnswer := fmt.Sprintf("@%s %s %d-%d %s", user, rikishi1.ShikonaEn, stat.RikishiWins, stat.OpponentWins, rikishi2.ShikonaEn)
	if len(stat.Matches) > 3 {
		wins := 0
		losses := 0
		count := 0
		for i := range stat.Matches {
			if count == 3 {
				break
			}
			if stat.Matches[i].WinnerID == 0 {
				continue
			}
			if stat.Matches[i].WinnerID == rikishi1.ID {
				wins++
			} else {
				losses++
			}
			count++
		}
		if count == 3 {
			finalAnswer += fmt.Sprintf(" (%d-%d)", wins, losses)
		}
	}
	return finalAnswer
}

func getLastMatches(name string, user string, httpClient *http.Client, rikishisList []*rikishi, r response) string {
	rikishi := getbyShikonaEn(name, rikishisList)
	if rikishi == nil {
		return fmt.Sprintf("@%s %s", user, r.ErrNotFoundShikona)
	}
	id := rikishi.ID
	url := fmt.Sprintf("https://www.sumo-api.com/api/rikishi/%d/matches", id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
	}
	limitBody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data matchUtil
	err = json.NewDecoder(limitBody).Decode(&data)
	if err != nil {
		return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
	}
	if len(data.Records) == 0 {
		return fmt.Sprintf("@%s %s", user, r.ErrNoRecentMatches)
	}
	var (
		opponent string
		res      string
	)
	answer := fmt.Sprintf("@%s %s: ", user, rikishi.ShikonaEn)
	m := 0
	for _, d := range data.Records {
		if m >= 3 {
			break
		}
		if d.WinnerID == 0 {
			continue
		}
		if d.EastID == id {
			opponent = d.WestShikona
		} else {
			opponent = d.EastShikona
		}
		if d.WinnerID == id {
			res = "[W]"
		} else {
			res = "[L]"
		}
		if m > 0 {
			answer += " | "
		}
		answer += fmt.Sprintf("%s %s (%s)", res, opponent, d.Kimarite)
		m++
	}
	return answer
}

func (b *bot) getNextMatch(name string, user string, r response) string {
	rikishi := getbyShikonaEn(name, b.rikishiList)
	if rikishi == nil {
		return fmt.Sprintf("@%s %s", user, r.ErrNotFoundShikona)
	}
	day, check := b.getBashoDayById(r)

	if check != "" {
		switch check {
		case "not ok":
			return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
		case "even month":
			return fmt.Sprintf("@%s %s", user, r.ErrNoBashoMonth)
		default:
			return fmt.Sprintf("@%s %s", user, check)
		}
	}
	torikumi, check := getTorikumiByDayId(b.currBasho.Date, day, b.httpClient, r)
	if check != "" {
		return fmt.Sprintf("@%s %s", user, check)
	}
	var opp string
	var oppId int
	var match match
	for i, m := range torikumi.Torikumi {
		if m.EastShikona == rikishi.ShikonaEn {
			opp = m.WestShikona
			oppId = m.WestID
			match = torikumi.Torikumi[i]
			break
		} else if m.WestShikona == rikishi.ShikonaEn {
			opp = m.EastShikona
			oppId = m.EastID
			match = torikumi.Torikumi[i]
			break
		}
	}
	if opp == "" {
		return fmt.Sprintf("@%s %s", user, r.ErrNoBoutMakuuchi)
	}
	var finalAnswer string
	var res string
	if match.WinnerID != 0 {
		if match.WinnerID == rikishi.ID {
			res = "[W]"
		} else {
			res = "[L]"
		}
		finalAnswer = fmt.Sprintf("@%s %s: %s %s (%s)", user, rikishi.ShikonaEn, res, opp, match.Kimarite)
	} else {
		mu, err := getMatchupById(rikishi.ID, oppId, b.httpClient)

		if err != nil {
			finalAnswer = fmt.Sprintf("@%s %s vs %s (%s #%d)", user, rikishi.ShikonaEn, opp, r.Match, match.MatchNo)
		} else {
			finalAnswer = fmt.Sprintf("@%s %s (%d) vs (%d) %s (%s #%d)", user, rikishi.ShikonaEn, mu.RikishiWins, mu.OpponentWins, opp, r.Match, match.MatchNo)
		}
	}
	return finalAnswer
}

func (b *bot) getTop5Bouts(user string, r response) string {
	day, check := b.getBashoDayById(r)

	if check != "" {
		switch check {
		case "not ok":
			return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
		case "even month":
			return fmt.Sprintf("@%s %s", user, r.ErrNoBashoMonth)
		default:
			return fmt.Sprintf("@%s %s", user, check)
		}
	}
	torikumi, check := getTorikumiByDayId(b.currBasho.Date, day, b.httpClient, r)
	if check != "" {
		return fmt.Sprintf("@%s %s", user, check)
	}
	finalAnswer := fmt.Sprintf("@%s [%s %d] ", user, r.Day, torikumi.Torikumi[0].Day)
	limit := len(torikumi.Torikumi)
	start := limit - 5
	if start < 0 {
		start = 0
	}

	for i := start; i < limit; i++ {
		match := torikumi.Torikumi[i]
		switch {
		case match.WinnerID == 0:
			finalAnswer += fmt.Sprintf("%s vs %s", match.EastShikona, match.WestShikona)
		case match.WinnerID == match.EastID:
			finalAnswer += fmt.Sprintf("%s [W] vs %s (%s)", match.EastShikona, match.WestShikona, match.Kimarite)
		case match.WinnerID == match.WestID:
			finalAnswer += fmt.Sprintf("%s vs [W] %s (%s)", match.EastShikona, match.WestShikona, match.Kimarite)
		}
		if i != limit-1 {
			finalAnswer += " | "
		}
	}
	return finalAnswer
}

func (b *bot) getCurrentScore(name string, user string, r response) string {
	rikishi := getbyShikonaEn(name, b.rikishiList)
	if rikishi == nil {
		return fmt.Sprintf("@%s %s", user, r.ErrNotFoundShikona)
	}
	_, check := b.getBashoDayById(r)

	if check != "" {
		switch check {
		case "not ok":
			return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
		case "even month":
			return fmt.Sprintf("@%s %s", user, r.ErrNoBashoMonth)
		default:
			return fmt.Sprintf("@%s %s", user, check)
		}
	}

	bashoId := b.currBasho.Date
	id := rikishi.ID

	url := fmt.Sprintf("https://www.sumo-api.com/api/rikishi/%d/matches?bashoId=%s", id, bashoId)
	resp, err := b.httpClient.Get(url)
	if err != nil {
		return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
	}
	limitBody := http.MaxBytesReader(nil, resp.Body, 5<<20)
	var data matchUtil
	err = json.NewDecoder(limitBody).Decode(&data)
	if err != nil {
		return fmt.Sprintf("@%s %s", user, r.ErrTechnical)
	}

	if len(data.Records) == 0 {
		return fmt.Sprintf("@%s %s %s", user, rikishi.ShikonaEn, r.ErrNoBoutsThisBasho)
	}
	var (
		wins   int
		losses int
	)
	for _, match := range data.Records {
		if match.WinnerID == 0 {
			continue
		}
		if match.WinnerID == id {
			wins++
		} else {
			losses++
		}
	}
	return fmt.Sprintf("@%s %s %d-%d", user, rikishi.ShikonaEn, wins, losses)
}
