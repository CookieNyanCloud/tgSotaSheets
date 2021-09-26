package main

import (
	"context"
	"errors"
	"github.com/cookienyancloud/tgSotaSheets/configs"
	"github.com/cookienyancloud/tgSotaSheets/internal/delivery"
	"github.com/cookienyancloud/tgSotaSheets/sotatgbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	ginServer "github.com/cookienyancloud/tgSotaSheets/internal/server"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	tokenA = "TOKEN_A"
	tokenB = "TOKEN_B"
	credFile = "driveapisearch.json"

)
const (
	pull = "1"
	push = "2"
)

const(
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
	bot, updates := sotatgbot.StartSotaBot(tokenB)

	for update := range updates {

		curUser := update.Message.ForwardFrom.UserName
		if update.Message.IsCommand() {
			println("command")

			_, ok := users[curUser]
			if !ok {
				println("no")
				continue
			}

			println("yes")
			users[curUser] = update.Message.Command()
			switch users[curUser] {
			case pull:
				msg1:= tgbotapi.NewMessage(update.Message.Chat.ID,pullWelcome)
				msg1.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg1)
			case push:

			default:
				msg1:= tgbotapi.NewMessage(update.Message.Chat.ID,pullWelcome)
				msg1.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg1)
			}

			continue
		}

		com:= users[curUser]
		switch com {
		case pull:

		}


	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

}
