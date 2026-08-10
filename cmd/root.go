package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ekucher/kbot/internal/handlers"
	"github.com/spf13/cobra"
	tele "gopkg.in/telebot.v4"
)

var rootCmd = &cobra.Command{
	Use:   "telegram-bot",
	Short: "Telegram bot written in Go",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := os.Getenv("TELE_TOKEN")

		if token == "" {
			return fmt.Errorf("environment variable TELE_TOKEN is not set")
		}

		bot, err := tele.NewBot(tele.Settings{
			Token: token,
			Poller: &tele.LongPoller{
				Timeout: 10 * time.Second,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create Telegram bot: %w", err)
		}

		handlers.Register(bot)

		fmt.Println("Telegram bot started")
		fmt.Println("Press Ctrl+C to stop")

		bot.Start()
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
