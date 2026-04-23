package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
)

func capitalizeFirstLetter(str string) string {
	runes := []rune(str)
	if len(runes) == 0 {
		return ""
	}
	return strings.ToUpper(string(runes[0])) + strings.ToLower(string(runes[1:]))
}

func translitShikona(name string, translitMap map[rune]string) string {
	name = strings.ToLower(name)
	runes := []rune(name)
	var shikona string

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if i+1 < len(runes) && r == 'д' && (runes[i+1] == 'ж' || runes[i+1] == 'з') {
			continue
		}
		if l, ok := translitMap[r]; ok {
			shikona += l
		} else {
			shikona += string(r)
		}
	}

	return shikona
}

func getbyShikonaEn(name string, rikishisList []*rikishi) *rikishi {
	name = capitalizeFirstLetter(name)
	var bestMatch *rikishi
	var limit int
	switch {
	case len(name) < 4:
		limit = 2
	case len(name) < 6:
		limit = 3
	default:
		limit = 4
	}

	bestDistance := limit
	for _, rikishi := range rikishisList {
		dist := levenshtein.ComputeDistance(name, rikishi.ShikonaEn)
		if dist == 0 {
			return rikishi
		}
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = rikishi
		}
	}
	return bestMatch
}

func getBashoId() string {
	now := time.Now().UTC()
	currMonth := int(now.Month())
	currYear := now.Year()
	id := fmt.Sprintf("%d%02d", currYear, currMonth)
	return id
}

func calculateDay(b *bashoTime, r response) (int, string) {
	now := time.Now().UTC()

	if now.Before(b.StartDate) {
		return 0, r.ErrBashoDontStart
	}

	day := int(now.Sub(b.StartDate).Hours()/24) + 1

	if now.Hour() > 9 || (now.Hour() == 9 && now.Minute() > 15) {
		day++
	}
	if day >= 16 {
		return 0, r.ErrBashoEnd
	}
	return day, ""
}
