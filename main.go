package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

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
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(port, nil)
	fmt.Println("start listen :", port)

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60

	// updates, err := bot.GetUpdatesChan(u)

	// получаем все обновления из канала updates
	for update := range updates {
		fmt.Println("update", update.Message.From, update.Message.From.ID, update.Message.From.LastName)
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "type /start or /settings."
		case "start":
			msg.Text = "Hi :)"
		case "add":
			fmt.Println("add", update.Message.CommandArguments())
			secret := &Secret{
				userId:  update.Message.From.ID,
				content: update.Message.CommandArguments(),
			}
			id, err := secret.save()
			if err != nil {
				msg.Text = "Error. Your content was not added"
				continue
			}
			msg.Text = "Your content was added " + strconv.Itoa(id)
		case "settings":
			msg.Text = "I know nothing about settings"
		default:
			msg.Text = "I don't know that command"
		}
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}
