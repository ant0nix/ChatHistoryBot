package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ant0nix/ChatHistoryBot/pkg/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

func (b *Bot) Run(db *sqlx.DB) error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates := b.initUPDChanel()
	b.handlerUPD(&updates, db)
	return nil
}

func (b *Bot) handlerUPD(updates *tgbotapi.UpdatesChannel, db *sqlx.DB) {

	mess_id := 0
	for update := range *updates {
		db := repository.NewDataBase(db)
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			db.NewRowIntoDB(&update)
			if update.Message.MessageID == mess_id {
				mess_id = -1
				err := db.NewAdmin(update)
				if err != nil {
					errorHadler(err, "(newAdmin)")
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Был добавлен новый администратор")
				_, err = b.bot.Send(msg)
				if err != nil {
					errorHadler(err, "(newAdmin/Sending)")
				}

			}
			if update.Message.Chat.IsPrivate() {
				switch update.Message.Text {
				case "/start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, repository.StartMessage)
					msg.ParseMode = "html"
					_, err := b.bot.Send(msg)
					errorHadler(err, "(start)")
				case "/show_all_users":
					answer, err := db.ShowAllUsers(update)
					errorHadler(err, "(ShowAllUsers)")
					err = b.SendingMessage(answer.Text, update)
					errorHadler(err, "(ShowAllUsers/Sending)")
				case "/show_all_messages":
					answer, err := db.ShowAllMessages(update)
					errorHadler(err, "(ShowAllMessages)")
					err = b.SendingMessage(answer.Text, update)
					errorHadler(err, "(ShowAllMessages/Sending)")
				case "/show_users_messages":
					users, err := db.ReturnUsers(update)
					if err != nil && err.Error() == "no admin" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You aren't admin of this bot!")
						b.bot.Send(msg)
						break
					}
					if err != nil && err.Error() != "no admin" {
						errorHadler(err, "(ShowUsersMess/ReturnUsers)")
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите внизу интересующего вас пользователя.")
					keyboard := tgbotapi.InlineKeyboardMarkup{}
					for _, values := range users {
						var row []tgbotapi.InlineKeyboardButton
						userid := strconv.Itoa(values.UserID)
						btn := tgbotapi.NewInlineKeyboardButtonData(values.UserName, "user_"+userid)
						row = append(row, btn)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}
					msg.ReplyMarkup = keyboard
					_, err = b.bot.Send(msg)
					errorHadler(err, "(ShowUsersMess/Send)")
				case "/show_chats_messages":
					users, err := db.ReturnMessages(update)
					if err != nil && err.Error() == "no admin" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You aren't admin of this bot!")
						b.bot.Send(msg)
						break
					}
					if err != nil && err.Error() != "no admin" {
						errorHadler(err, "(ShowChatsMess/ReturnUsers)")
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите внизу интересующий вас чат.")
					keyboard := tgbotapi.InlineKeyboardMarkup{}
					for _, values := range users {
						var row []tgbotapi.InlineKeyboardButton
						chatid := strconv.Itoa(values.ChatID)
						btn := tgbotapi.NewInlineKeyboardButtonData(values.ChatTitle, "chat_"+chatid)
						row = append(row, btn)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}
					msg.ReplyMarkup = keyboard
					_, err = b.bot.Send(msg)
					errorHadler(err, "(ShowChatMess/Send)")
				case "/newAdmin":
					mess_id = update.Message.MessageID + 2
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пришлите сообщение человека, которого нужно назначит админом")
					b.bot.Send(msg)
				}
			}
		}
		if update.CallbackQuery != nil {

			callback := strings.Split(update.CallbackQuery.Data, "_")
			fmt.Println(callback)
			switch callback[0] {
			case "user":
				answer, err := db.SelectOneUser(update)
				errorHadler(err, "(Callback/SelectOneUser)")
				err = b.SendingMessageCallbak(answer.Text, update)
				errorHadler(err, "(Callback/SelectOneUser/Send)")
			case "chat":
				answer, err := db.SelectOneChat(update)
				errorHadler(err, "(Callback/SelectOneChat)")
				err = b.SendingMessageCallbak(answer.Text, update)
				errorHadler(err, "(Callback/SelectOneChat/Send)")
			}
		}
	}
}

func (b *Bot) initUPDChanel() tgbotapi.UpdatesChannel {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) SendingMessage(messages string, update tgbotapi.Update) error {
	fmt.Println(len(messages))
	for i := 0; i < len(messages); {
		fmt.Println(len(messages))
		if len(messages) < 4000 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages)
			_, err := b.bot.Send(msg)
			return err
		} else if i+4000 < len(messages) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages[i:i+4000])
			i += 4000
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			fmt.Println(len(messages))
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages[i:])
			_, err := b.bot.Send(msg)
			return err
		}
	}
	return nil
}

func (b *Bot) SendingMessageCallbak(messages string, update tgbotapi.Update) error {
	fmt.Println(len(messages))
	for i := 0; i < len(messages); {
		fmt.Println(len(messages))
		if len(messages) < 4000 {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages)
			_, err := b.bot.Send(msg)
			return err
		} else if i+4000 < len(messages) {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages[i:i+4000])
			i += 4000
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			fmt.Println(len(messages))
		} else {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages[i:])
			_, err := b.bot.Send(msg)
			return err
		}
	}
	return nil
}

func errorHadler(err error, funcName string) {
	if err != nil {
		log.Printf("ERROR!In func %s. Error text: %s", funcName, err.Error())
	}
}
