package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func handler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	log.Println("update", msg.From, msg.From.ID, msg.From.LastName)
	if msg == nil { // ignore any non-Message Updates
		return
	}

	if !msg.IsCommand() { // ignore any non-command Messages
		return
	}

	var text string
	switch msg.Command() {
	case "help":
		text = "type /start or /settings."
	case "start":
		text = "Hi :)"
	case "add":
		log.Println("add", msg.CommandArguments())
		secret := &Secret{
			userId:  msg.From.ID,
			content: msg.CommandArguments(),
		}
		id, err := secret.save()
		if err != nil {
			text = "Error. Your content was not added"
			return
		}
		text = "Your content was added " + strconv.Itoa(id)
	case "settings":
		text = "I know nothing about settings"
	default:
		text = "I don't know that command"
	}

	newMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
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
		go handler(bot, update.Message)
	}
}
