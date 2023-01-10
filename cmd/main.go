package main

import (
	"log"
	"os"

	"github.com/ant0nix/ChatHistoryBot/pkg/repository"
	"github.com/ant0nix/ChatHistoryBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func main() {
	if err := gotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err)
	}

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err)
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("filed to inizialize db: %s\n", err.Error())
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	tBot := telegram.NewBot(bot)
	if err := tBot.Run(db); err != nil {
		log.Fatal()
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
