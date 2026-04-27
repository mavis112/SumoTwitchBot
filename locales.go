package main

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "w", 'г': "g", 'д': "d",
	'е': "e", 'ё': "yo", 'ж': "j", 'з': "z", 'и': "i",
	'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
	'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
	'у': "u", 'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch",
	'ш': "sh", 'щ': "sh", 'ы': "i", 'э': "e", 'ю': "yu",
	'я': "ya",
}

type response struct {
	ErrNoRikishi       string
	ErrNotFoundShikona string
	ErrBadFormatDebut  string
	CurrentRank        string
	Heya               string
	Height             string
	Weight             string
	Age                string
	Debut              string
	HighestRank        string
	Matches            string
	Match              string
	Wins               string
	NumOfBasho         string
	Yusho              string
	Makuuchi           string
	Juryo              string
	Day                string
	ErrNoMatchup       string
	ErrNoRecentMatches string
	ErrTechnical       string
	ErrNoBashoMonth    string
	ErrBashoDontStart  string
	ErrBashoEnd        string
	ErrNoScheduleYet   string
	ErrNoBoutMakuuchi  string
	ModOnlyOn          string
	ModOnlyOff         string
	NeedTranslit       bool
}

var ruResp = response{
	ErrNoRikishi:       "укажи имена борцов",
	ErrNotFoundShikona: "не удалось найти борца с таким именем",
	ErrBadFormatDebut:  "н/д",
	CurrentRank:        "Ранг",
	Heya:               "Школа",
	Height:             "см",
	Weight:             "кг",
	Age:                "Возраст",
	Debut:              "Дебют",
	HighestRank:        "Макс. ранг",
	Matches:            "Матчи",
	Match:              "матч",
	Wins:               "Победы",
	NumOfBasho:         "Турниры",
	Yusho:              "Кубки",
	Makuuchi:           "Макуучи",
	Juryo:              "Джурё",
	Day:                "День",
	ErrNoMatchup:       "борцы ещё не встречались на дохё",
	ErrNoRecentMatches: "данные о матчах не найдены",
	ErrTechnical:       "не удалось получить данные",
	ErrNoBashoMonth:    "в этом месяце нет турнира",
	ErrBashoDontStart:  "турнир еще не начался",
	ErrBashoEnd:        "турнир уже завершен",
	ErrNoScheduleYet:   "расписание на этот день еще не опубликовано",
	ErrNoBoutMakuuchi:  "борец сегодня не выступает в Макуути",
	ModOnlyOn:          `включен режим "только для модераторов"`,
	ModOnlyOff:         `режим "только для модераторов" выключен. Команды доступны всем.`,
	NeedTranslit:       true,
}

var enResp = response{
	ErrNoRikishi:       "please specify rikishi names",
	ErrNotFoundShikona: "rikishi not found",
	ErrBadFormatDebut:  "N/A",
	CurrentRank:        "Rank",
	Heya:               "Heya",
	Height:             "cm",
	Weight:             "kg",
	Age:                "Age",
	Debut:              "Debut",
	HighestRank:        "Peak rank",
	Matches:            "Bouts",
	Match:              "bout",
	Wins:               "Wins",
	NumOfBasho:         "Basho",
	Yusho:              "Yusho",
	Makuuchi:           "Makuuchi",
	Juryo:              "Juryo",
	Day:                "Day",
	ErrNoMatchup:       "no head-to-head matches found",
	ErrNoRecentMatches: "no recent bouts found",
	ErrTechnical:       "failed to get data",
	ErrNoBashoMonth:    "no basho is scheduled for this month",
	ErrBashoDontStart:  "the tournament hasn't started yet",
	ErrBashoEnd:        "the tournament has already concluded",
	ErrNoScheduleYet:   "schedule for this day is not available yet",
	ErrNoBoutMakuuchi:  "the rikishi is not competing in Makuuchi today",
	ModOnlyOn:          `"moderators only" mode activated.`,
	ModOnlyOff:         `"moderators only" mode deactivated. All users can now use commands`,
	NeedTranslit:       false,
}
