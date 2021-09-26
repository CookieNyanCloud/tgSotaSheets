package main

import (
	"context"
	"fmt"
	"github.com/cookienyancloud/tgSotaSheets/configs"
	"github.com/cookienyancloud/tgSotaSheets/internal/service"
	"github.com/cookienyancloud/tgSotaSheets/sotatgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	tokenA   = "TOKEN_A"
	tokenB   = "TOKEN_B"
	credFile = "driveapisearch.json"
)
const (
	pull = "1"
	push = "2"
)

const (
	pullWelcome = "пришлите значение поиска"
	pushWelcome = `пришлите контакт в формате "ФИО, должность, номер, тг, дополнительно"
					пустыне строки можно пропустить`
)

func main() {

	conf := configs.InitConf()
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatalf("Unable to parse credantials file: %v", err)
	}

	users := configs.InitUsers()
	bot, updates := sotatgbot.StartSotaBot(tokenA)

	for update := range updates {
		if update.Message.Command() == "start" {
			continue
		}
		curUser := update.Message.From.UserName
		if update.Message.IsCommand() {
			_, ok := users[curUser]
			if !ok {
				continue
			}

			users[curUser] = update.Message.Command()
			switch users[curUser] {
			case pull:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, pullWelcome)
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			case push:

			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, pullWelcome)
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}

			continue
		}
		if update.Message.Text == ""{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "не текст")
			_, _ = bot.Send(msg)
			continue
		}
		com, ok := users[curUser]
		if !ok {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, pullWelcome)
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
			continue
		}
		switch com {
		case pull:

			res, err := service.GetContact(srv, conf.SheetsAdr, update.Message.Text)
			if err != nil {
				fmt.Println("in case:", err)
				continue
			}
			var message string
			var format string
			for _= range res {
				format+=`%s`
			}
			for i, contact := range res {
				message += fmt.Sprintf("%v)%v,%v,%v,%v,%v,\n", i+1, contact.Name, contact.Job, contact.Cell, contact.Tg, contact.Other)
			}
			if message==""{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "в базе нет")
				_, _ = bot.Send(msg)
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			_, _ = bot.Send(msg)
			continue
		}
		continue
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

}
