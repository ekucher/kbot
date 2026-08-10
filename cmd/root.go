package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbot",
	Short: "Telegram bot written in Go",
	Long:  "Telegram bot written in Go using Cobra and Telebot.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
