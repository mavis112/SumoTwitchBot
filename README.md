# SumoTwitchBot

> ⚠️ **Примечание:** Готовый к запуску файл (`.exe`) будет опубликован в разделе Releases в промежутке между объявлением Бандзукэ (27 апреля) и началом Нацу Басё (10 мая). Сейчас в репозитории доступен исходный код для самостоятельной сборки.


> ⚠️ **Note:** The ready-to-run binary (`.exe`) will be published in the Releases section between the Banzuke announcement (April 27) and the start of the Natsu Basho (May 10). Currently, only the source code is available for manual building.


[RU (Русский)](#russian-версия) | [EN (English)](#english-version)

---

## Russian Версия

Легкий и быстрый Twitch-бот для фанатов сумо. Позволяет зрителям узнавать статистику рикиши, историю встреч и информацию о предстоящих поединках прямо в чате. 

Бот работает на данных [Sumo-API](https://sumo-api.com). При запуске бот автоматически загружает базу из 600+ активных борцов. Если загрузка не удалась (проблемы с сетью/API) — бот не запустится.

### Особенности
*   **Умный поиск:** Находит борцов даже с опечатками. 
    *   *Важно:* Старайтесь не сокращать имена слишком сильно (например, запрос `Хощ` скорее всего найдет `Oho`, а не `Hoshoryu`).
*   **Гибкость:** Понимает кириллицу (Поливанов + общепринятый транслит). Бот нечувствителен к регистру и прочитает команду, даже если после неё в сообщении есть другой текст.
*   **Безопасность:** Бот работает только через IRC-протокол. Не требует Client ID и не запрашивает лишних прав.

### Установка и запуск
1. Скачайте последнюю версию из раздела **Releases**.
2. Распакуйте архив и найдите файл `config.env`.
3. Заполните данные через Блокнот:
   * `CHANNEL`: Канал, куда должен зайти бот.
   * `BOT_NAME`: Имя аккаунта, от лица которого бот будет писать в чат. Если у вас нет отдельного аккаунта для бота, просто впишите имя своего основного канала.
   * `OAUTH_TOKEN`: Код доступа (Access Token). Важно: токен должен быть сформирован именно для того аккаунта, который указан в BOT_NAME.


#### Как получить токен:
1. Перейдите на сайт [TwitchTokenGenerator](https://twitchtokengenerator.com).
2. Нажмите на кнопку **Chat Bot**.
3. **Важно:** В списке разрешений (Scopes) оставьте галочки **ТОЛЬКО** на этих двух пунктах:
   * `chat:read` — чтение команд в чате.
   * `chat:edit` — отправка ответов.
   * **Все остальные галочки нужно СНЯТЬ.**
4. Нажмите «Generate Token», авторизуйтесь и скопируйте полученный **Access Token**. 

### Пример заполнения config.env:

```env
CHANNEL=twitch_user69
BOT_NAME=sumo_bot_active
OAUTH_TOKEN=u7823h4iu23hi4u23hi4u23h
```

#### Как обновлять бота:
При выходе новой версии достаточно заменить старый файл `sumobot.exe` на новый. **Не перезаписывайте** уже настроенный файл `config.env`, чтобы не потерять свои данные.

### Команды


| Команда | Описание | Пример |
| :--- | :--- | :--- |
| `!стат [имя]` | Ранг, хэя, рост и вес борца и др. | *!стат рога* |
| `!ласт [имя]` | Результаты 3-х последних матчей. | *!ласт шиши* |
| `!матчап [имя1] [имя2]` | История личных встреч борцов. | *!матчап аби оносато* |
| `!след [имя]` | Соперник и статистика матча **(только дивизион Макуучи)**. | *!след аонишики* |
| `!топ5` | 5 последних поединков игрового дня **(только дивизион Макуучи)**. | *!топ5* |

> **Примечание:** Если борец уже выступил сегодня, `!след` покажет результат матча. Если турнир сейчас не идет, команда выдаст ошибку.

---

## English Version

A lightweight and fast Twitch bot for sumo fans. It allows viewers to check rikishi stats, head-to-head history, and upcoming match information directly in chat.

The bot uses data from [Sumo-API](https://sumo-api.com). On startup, it automatically downloads a database of 600+ active wrestlers.

### Features
*   **Smart Search:** Finds wrestlers even with typos (Levenshtein distance).
    *   *Note:* Avoid very short abbreviations (e.g., `Hosh` might return `Oho` instead of `Hoshoryu`).
*   **Multilingual:** Supports both English names and Cyrillic transliteration.
*   **Flexible Parsing:** Case-insensitive and ignores extra text after the command.
*   **Safety:** Uses only the IRC protocol; no Client ID or unnecessary permissions required.

### Setup & Launch
1. Download the latest release from the **Releases** section.
2. Unzip and open the `config.env` file with Notepad.
3. Fill in the following:
   * `CHANNEL`: The channel the bot will join.
   * `BOT_NAME`: The name of the account that will post in chat. If you don't have a dedicated bot account, just use your main channel name.
   * `OAUTH_TOKEN`: Access Token. Important: the token must be generated for the exact account specified in BOT_NAME.


### config.env Example:

```env
CHANNEL=twitch_user69
BOT_NAME=sumo_bot_active
OAUTH_TOKEN=u7823h4iu23hi4u23hi4u23h
```

#### How to get a Token:
1. Go to [TwitchTokenGenerator](https://twitchtokengenerator.com).
2. Select **Chat Bot**.
3. **Important:** In the Scopes list, check **ONLY** these two:
   * `chat:read`
   * `chat:edit`
   * **Uncheck everything else.**
4. Copy the **Access Token**.

#### How to Update:
When a new version is released, only replace the old `sumobot.exe` with the new one. **Do not overwrite** your existing `config.env` file.

### Commands


| Command | Description | Example |
| :--- | :--- | :--- |
| `!stats [name]` | Rank, heya, height, weight, and more | *!stats asanoyama* |
| `!last [name]` | Results of the last 3 matches. | *!last hoshoryu* |
| `!matchup [n1] [n2]` | Head-to-head history. | *!matchup abi onosato* |
| `!next [name]` | Next opponent & match stats **(Makuuchi only)**. | *!next kotozakura* |
| `!top5` | The last 5 bouts of the tournament day **(Makuuchi only)**. | *!top5* |

---

## ⚠️ Disclaimer / Статус проекта
Бот находится в режиме интенсивного тестирования. Критичные баги будут исправляться во время Басё. Распространяется по лицензии **MIT**.
The bot is in intensive testing mode. Critical bugs will be patched during Basho. Distributed under the **MIT License**.
