package main

import (
	"context"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"ttv-bot/config"
	"ttv-bot/service"
)

var gcron *cron.Cron

func main() {
	var cmd = cobra.Command{
		Use: "ttv-bot",
	}
	cmd.PersistentFlags().String(config.BotToken, "", "Bot token")
	cmd.PersistentFlags().Bool(config.BotDebugMode, false, "Bot Debug Mode")
	cmd.PersistentFlags().Int(config.BotTimeout, 60, "Bot Update Timeout")
	cmd.PersistentFlags().String(config.ApiKey, "demo", "Chainbase Api Key")
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Fatal("failed to bind flags", zap.Error(err))
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		gcron = cron.New()
		gcron.Start()

		botService := service.NewService(viper.GetString(config.BotToken), viper.GetBool(config.BotDebugMode), viper.GetInt(config.BotTimeout), viper.GetString(config.ApiKey))
		botService.Start()
		return nil
	}

	if err := cmd.ExecuteContext(context.Background()); err != nil {
		log.Fatal("failed to execute command", zap.Error(err))
	}

}
