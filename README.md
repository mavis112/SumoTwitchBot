# SumoTwitchBot

[RU (Русский)](#русская-версия) | [EN (English)](#english-version)

---

## Русская Версия

Лёгкий Twitch-бот для фанатов сумо. Позволяет зрителям узнавать статистику рикиши, историю встреч и информацию о матчах прямо в чате. 

Работает на данных [Sumo-API](https://sumo-api.com). При запуске загружает базу из 600+ активных борцов. Если загрузка не удалась (проблемы с сетью/API) — бот не запустится.

### Особенности
*   **Умный поиск:** Находит борцов даже с опечатками (Levenshtein distance). 
    *   *Важно:* Старайтесь не сокращать имена слишком сильно (например, запрос `Хощ` скорее всего найдет `Oho`, а не `Hoshoryu`).
*   **Гибкость:** Понимает кириллицу. Нечувствителен к регистру и лишнему тексту после команды.
*   **Безопасность:** Только IRC-протокол. Не требует Client ID.

### Установка и запуск
1. Скачайте последнюю версию из раздела **Releases**.
2. Распакуйте архив и найдите файл `config.env`.
3. Заполните данные через Блокнот:
   *   `CHANNEL`: Канал, куда должен зайти бот.
   *   `BOT_NAME`: Имя аккаунта, от лица которого бот будет писать в чат. Если у вас нет отдельного аккаунта для бота, просто впишите имя своего основного канала.
   *   `OAUTH_TOKEN`: Код доступа (Access Token). **Важно:** токен должен быть сформирован именно для того аккаунта, который указан в `BOT_NAME`.

#### Как получить токен:
1. Перейдите на сайт [TwitchTokenGenerator](https://twitchtokengenerator.com).
2. Нажмите на кнопку **Chat Bot**.
3. **Важно:** В списке разрешений (Scopes) оставьте галочки **ТОЛЬКО** на этих двух пунктах:
   *   `chat:read` — чтение команд в чате.
   *   `chat:edit` — отправка ответов.
   *   **Все остальные галочки нужно СНЯТЬ.**
4. Нажмите «Generate Token», авторизуйтесь и скопируйте полученный **Access Token**. 

> **Рекомендация:** Если бот работает с отдельного аккаунта, дайте ему статус модератора (`/mod имя_бота`), чтобы избежать лимитов Twitch на отправку сообщений.

### Пример заполнения config.env:
```env
CHANNEL=twitch_user69
BOT_NAME=sumo_bot_active
OAUTH_TOKEN=u7823h4iu23hi4u23hi4u23h
```

#### Как обновлять бота:
При выходе новой версии достаточно заменить старый файл `SumoTwitchBot.exe` на новый. **Не перезаписывайте** уже настроенный файл `config.env`, чтобы не потерять свои данные.

### Защита от спама
*   **3 секунды** кулдауна на использование команд для обычных зрителей.
*   Стример и модераторы — **без ограничений**.

### Команды


| Команда | Описание | Пример |
| :--- | :--- | :--- |
| `!стат [имя]` | Ранг, хэя, рост и вес борца и др. | *!стат рога* |
| `!ласт [имя]` | Результаты 3-х последних матчей. | *!ласт шиши* |
| `!матчап [имя1] [имя2]` | История личных встреч борцов. | *!матчап аби оносато* |
| `!след [имя]` | Соперник и статистика матча **(только дивизион Макуучи)**. | *!след аонишики* |
| `!топ5` | 5 последних поединков игрового дня **(только дивизион Макуучи)**. | *!топ5* |
| `!модрежим вкл/выкл` | Режима «только для модераторов». [Стример/Модераторы] | *!модрежим вкл* |

---

## English Version

A lightweight Twitch bot for sumo fans. It allows viewers to check rikishi stats, head-to-head history, and match info directly in chat.

Uses data from [Sumo-API](https://sumo-api.com). On startup, it downloads a database of 600+ active wrestlers. If the download fails (network/API issues), the bot will not start.

### Features
*   **Smart Search:** Finds wrestlers even with typos (Levenshtein distance).
    *   *Note:* Avoid very short abbreviations (e.g., `Hosh` might return `Oho` instead of `Hoshoryu`).
*   **Multilingual:** Supports both English names and Cyrillic transliteration.
*   **Safety:** Uses IRC protocol only. No Client ID required.

### Setup & Launch
1. Download the latest version from the **Releases** section.
2. Unzip the archive and find the `config.env` file.
3. Fill in the data using Notepad:
   *   `CHANNEL`: The channel the bot should join.
   *   `BOT_NAME`: The account name the bot will post as. If you don't have a separate bot account, just use your main channel name.
   *   `OAUTH_TOKEN`: Access Token. **Important:** the token must be generated for the exact account specified in `BOT_NAME`.

#### How to get a token:
1. Go to [TwitchTokenGenerator](https://twitchtokengenerator.com).
2. Click the **Chat Bot** button.
3. **Important:** In the Scopes list, leave **ONLY** these two checked:
   *   `chat:read`
   *   `chat:edit`
   *   **UNCHECK everything else.**
4. Click "Generate Token", authorize, and copy the **Access Token**.

> **Recommendation:** If using a dedicated bot account, grant it moderator status (`/mod bot_name`) to avoid Twitch message rate limits.

### config.env Example:
```env
CHANNEL=twitch_user69
BOT_NAME=sumo_bot_active
OAUTH_TOKEN=u7823h4iu23hi4u23hi4u23h
```

#### How to update:
When a new version is released, simply replace the old `SumoTwitchBot.exe` with the new one. **DO NOT overwrite** your existing `config.env` file.

### Spam Protection
*   **3-second** command cooldown for regular viewers.
*   Broadcaster and moderators have **no limits**.

### Commands


| Command | Description | Example |
| :--- | :--- | :--- |
| `!stats [name]` | Rank, heya, height, weight, etc. | *!stats asanoyama* |
| `!last [name]` | Results of the last 3 matches. | *!last hoshoryu* |
| `!matchup [n1] [n2]` | Head-to-head history. | *!matchup abi onosato* |
| `!next [name]` | Next opponent & match stats **(Makuuchi only)**. | *!next kotozakura* |
| `!top5` | The final 5 bouts of the tournament day **(Makuuchi only)**. | *!top5* |
| `!modsonly on/off` | Toggle "Moderators Only" mode. [Broadcaster/Mods] | *!modsonly on* |

---

## ⚠️ Disclaimer / Статус проекта
Бот находится в режиме интенсивного тестирования. Критичные баги будут исправляться во время Басё. Распространяется по лицензии **MIT**.

The bot is in intensive testing mode. Critical bugs will be patched during Basho. Distributed under the **MIT License**.