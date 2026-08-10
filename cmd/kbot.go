package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ekucher/kbot/internal/handlers"
	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v4"
)

var (
	// TeleToken contains the Telegram Bot API token.
	TeleToken = os.Getenv("TELE_TOKEN")
)

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Start Telegram bot",
	Long:    "Start the Telegram bot using the TELE_TOKEN environment variable.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if TeleToken == "" {
			return fmt.Errorf("TELE_TOKEN environment variable is not set")
		}

		kbot, err := telebot.NewBot(telebot.Settings{
			Token: TeleToken,
			Poller: &telebot.LongPoller{
				Timeout: 10 * time.Second,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create Telegram bot: %w", err)
		}

		registerHandlers(kbot)

		log.Println("Telegram bot started")
		log.Println("Press Ctrl+C to stop")

		kbot.Start()

		return nil
	},
}

// startHandler intentionally uses telebot.Context directly.
// It also delegates the actual message processing to the handlers package.
func startHandler(c telebot.Context) error {
	return handlers.Start(c)
}

func registerHandlers(kbot *telebot.Bot) {
	kbot.Handle("/start", startHandler)
	kbot.Handle("/help", handlers.Help)
	kbot.Handle("/hello", handlers.Hello)

	kbot.Handle(telebot.OnText, handlers.Text)
	kbot.Handle(telebot.OnPhoto, handlers.Photo)
	kbot.Handle(telebot.OnDocument, handlers.Document)
	kbot.Handle(telebot.OnSticker, handlers.Sticker)
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}
