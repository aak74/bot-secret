package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var store *Store

func handleCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) string {
	log.Println("handleCommand", msg.Command(), msg.CommandArguments())
	switch msg.Command() {
	case "help":
		return "type /add or /settings."
	case "add":
		id, err := save(msg)
		if err != nil {
			return "Error. Your content was not added"
		}
		return "Your content was added " + strconv.Itoa(id)
	case "list":
		res, err := list(msg)
		if err != nil {
			return "GetList error"
		}
		return "list " + strings.Join(res[:], ",")
	case "settings":
		return "I know nothing about settings"
	default:
		return "I don't know that command"
	}
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log.Println("update", msg.From, msg.From.ID, msg.From.LastName)
	if msg == nil { // ignore any non-Message Updates
		return
	}

	if !msg.IsCommand() { // ignore any non-command Messages
		return
	}

	newMsg := tgbotapi.NewMessage(
		msg.Chat.ID,
		handleCommand(bot, msg))
	if _, err := bot.Send(newMsg); err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	webhookURL := os.Getenv("WEBHOOK_URL")
	port := os.Getenv("PORT")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	store = NewBoltDB()
	// bot.Debug = true
	log.Println("Authorized on account", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(port, nil)
	log.Println("start listen", port)

	for update := range updates {
		go handleMessage(bot, update.Message)
	}
}
