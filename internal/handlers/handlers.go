package handlers

import (
	"strings"

	tele "gopkg.in/telebot.v4"
)

func Register(bot *tele.Bot) {
	bot.Handle("/start", Start)
	bot.Handle("/help", Help)
	bot.Handle("/hello", Hello)

	bot.Handle(tele.OnText, Text)
	bot.Handle(tele.OnPhoto, Photo)
	bot.Handle(tele.OnDocument, Document)
	bot.Handle(tele.OnSticker, Sticker)
}

func Start(c tele.Context) error {
	return c.Send(
		"Привіт! Я Telegram-бот, написаний мовою Go.\n\n" +
			"Напиши /help, щоб переглянути доступні команди.",
	)
}

func Help(c tele.Context) error {
	return c.Send(
		"Доступні команди:\n" +
			"/start - запуск бота\n" +
			"/help - довідка\n" +
			"/hello - привітання\n\n" +
			"Текстові команди:\n" +
			"hello\n" +
			"привіт\n" +
			"ping\n" +
			"golang",
	)
}

func Hello(c tele.Context) error {
	return c.Send("Hello!")
}

func Text(c tele.Context) error {
	text := strings.TrimSpace(c.Text())

	switch strings.ToLower(text) {
	case "hello", "hi", "привіт", "вітаю":
		return c.Send("Привіт! Радий тебе бачити.")

	case "ping":
		return c.Send("pong")

	case "golang", "go":
		return c.Send("Go — чудова мова для створення Telegram-ботів.")

	default:
		return c.Send("Ви написали: " + text)
	}
}

func Photo(c tele.Context) error {
	return c.Send("Я отримав фотографію.")
}

func Document(c tele.Context) error {
	return c.Send("Я отримав документ.")
}

func Sticker(c tele.Context) error {
	return c.Send("Я отримав стікер.")
}
