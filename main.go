package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cookienyancloud/tgSotaSheets/configs"
	"github.com/cookienyancloud/tgSotaSheets/service"
	"github.com/cookienyancloud/tgSotaSheets/sotatgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	tokenA   = "TOKEN_A"
	credFile = "driveapisearch.json"
)
const (
	pull     = "1"
	push     = "2"
	pushUser = "3"
)

const (
	welcome     = "бот поиска по инкунабуле контактов v2"
	pullWelcome = "пришлите значение поиска"
	pushWelcome = `пришлите контакт в формате "ФИО, должность, номер, тг, дополнительно"
,пустыне строки можно пробелом`
	pushUserWelcome = `пришлите тг ник человека`
	restricted      = "в доступе отказано"
	unknown         = "хз"
)

func main() {

	conf := configs.InitConf()
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatalf("Unable to parse credantials file: %v", err)
	}

	users, err := configs.InitUsers()
	if err != nil {
		log.Fatalf("error getting users: %v", err)
	}

	bot, updates := sotatgbot.StartSotaBot(conf.Token)
	for update := range updates {

		if update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcome)
			_, _ = bot.Send(msg)
			continue
		}

		curUser := update.Message.From.UserName

		if update.Message.IsCommand() {
			_, ok := users[curUser]
			if !ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, restricted)
				_, _ = bot.Send(msg)
				continue
			}
			users[curUser] = update.Message.Command()
			switch users[curUser] {
			case pull:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, pullWelcome)
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			case push:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, pushWelcome)
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			case pushUser:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, pushUserWelcome)
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, unknown)
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}
			continue
		}

		if update.Message.Text == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "не текст")
			_, _ = bot.Send(msg)
			continue
		}

		com, ok := users[curUser]
		if !ok {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, restricted)
			_, _ = bot.Send(msg)
			continue
		}

		switch com {
		case pull:
			if  len(update.Message.Text) <5 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "слишком широкая выборка")
				_, _ = bot.Send(msg)
				break
			}

			res, err := service.GetContact(srv, conf.SheetsAdr, update.Message.Text)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln("in case pull:", err))
				_, _ = bot.Send(msg)
				fmt.Println("in case pull:", err)
				break
			}
			var message string
			for i, contact := range res {
				message += fmt.Sprintf("%v)", i+1)
				message += fmt.Sprintf("%v, ", contact.Name)
				if contact.Job != "" {
					message += fmt.Sprintf("%v, ", contact.Job)
				}
				if contact.Cell != "" {
					message += fmt.Sprintf("%v, ", contact.Cell)
				}
				if contact.Tg != "" {
					message += fmt.Sprintf("%v, ", contact.Tg)
				}
				if contact.Other != "" {
					message += fmt.Sprintf("%v", contact.Other)
				}
				message += fmt.Sprintf("\n")
			}
			if message == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "в базе нет")
				_, _ = bot.Send(msg)
				break
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			_, _ = bot.Send(msg)
		case push:
			input := strings.Split(update.Message.Text, ",")
			fmt.Println(input)
			if len(input) != 5 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "не формат")
				_, _ = bot.Send(msg)
				break
			}
			userContact := service.Result{
				Name:  input[0],
				Job:   input[1],
				Cell:  input[2],
				Tg:    input[3],
				Other: input[4],
			}
			err := service.SendContact(srv, conf.SheetsAdr, userContact)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln("in case push:", err))
				_, _ = bot.Send(msg)
				fmt.Println("in case:", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "добавлен")
			_, _ = bot.Send(msg)
		case pushUser:
			name := strings.ReplaceAll(update.Message.Text, "@", "")
			err := configs.AddUser(users,update.Message.From.UserName, name)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln("in case pushUser:", err))
				_, _ = bot.Send(msg)
				fmt.Println("in case pushUser:", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "добавлен")
			_, _ = bot.Send(msg)
		}
		continue
	}
}
