package entities

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Admin struct {
	UserID      int    `db:"userid"`
	UserName    string `db:"username"`
	MasterAdmin bool   `db:"masteradmin"`
}

func NewAdmin(update tgbotapi.Update) *Admin {
	return &Admin{
		UserID:      int(update.Message.From.ID),
		UserName:    update.Message.From.UserName,
		MasterAdmin: false,
	}
}
func NewAdminFow(update tgbotapi.Update) *Admin {
	return &Admin{
		UserID:      int(update.Message.ForwardFrom.ID),
		UserName:    update.Message.ForwardFrom.UserName,
		MasterAdmin: false,
	}
}

type User struct {
	ID       int    `db:"id"`
	UserID   int    `db:"userid"`
	FName    string `db:"fname"`
	LName    string `db:"lname"`
	UserName string `db:"username"`
	Language string `db:"languages"`
	ChatId   int    `db:"chatid"`
}

func NewUser(update *tgbotapi.Update) *User {
	return &User{
		UserID:   int(update.Message.From.ID),
		FName:    update.Message.From.FirstName,
		LName:    update.Message.From.LastName,
		UserName: update.Message.From.UserName,
		Language: update.Message.From.LanguageCode,
		ChatId:   int(update.Message.Chat.ID),
	}
}

type Message struct {
	ID        int    `db:"id"`
	Date      string `db:"mdate"`
	Owner     int    `db:"mowner"`
	InputData string `db:"inputdata"`
	ChatID    int    `db:"chatid"`
	ChatTitle string `db:"chattitle"`
}

func NewMessage(update *tgbotapi.Update) *Message {
	unix_timestamp := update.Message.Date
	tm := time.Unix(int64(unix_timestamp), 0).Format("2 January 2006 15:04")
	tm2 := update.Message.Chat.Title
	if update.Message.Chat.Title == "" {
		tm2 = "privite_mess_" + update.Message.From.UserName
	}
	return &Message{
		Date:      tm,
		InputData: update.Message.Text,
		ChatID:    int(update.Message.Chat.ID),
		ChatTitle: tm2,
	}
}
