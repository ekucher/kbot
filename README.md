\# Telegram Bot in Go



Simple Telegram bot written in Go using Cobra and Telebot v4.



\## Technologies



\- Go

\- Cobra

\- Telebot v4

\- Telegram Bot API



\## Features



The bot supports:



\- `/start` - start the bot

\- `/help` - show available commands

\- `/hello` - greeting

\- text messages

\- photos

\- documents

\- stickers



The bot also processes text message content.



Examples:



```text

hello

привіт

ping

golang

```



\## Requirements



\- Go 1.26 or newer

\- Telegram Bot token



\## Installation



Clone the repository:



```bash

git clone https://github.com/ekucher/telegram-bot.git

cd telegram-bot

```



Install dependencies:



```bash

go mod download

```



\## Configuration



Create a Telegram bot using BotFather.



Set the Telegram token using the `TELE\_TOKEN` environment variable.



\### Windows PowerShell



```powershell

$env:TELE\_TOKEN = "YOUR\_BOT\_TOKEN"

```



\### Linux



```bash

export TELE\_TOKEN="YOUR\_BOT\_TOKEN"

```



Do not store the Telegram Bot token in the source code or repository.



\## Run



```bash

go run .

```



\## Build



```bash

go build ./...

```



\## Telegram Bot



https://t.me/ekucher_go_bot

