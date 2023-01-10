package repository

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/ant0nix/ChatHistoryBot/pkg/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap/buffer"
)

func CheckAccess(update tgbotapi.Update, db *sqlx.DB) bool {
	var admins []entities.Admin
	err := db.Select(&admins, "SELECT * FROM admins")
	if err != nil {
		log.Println(err)
	}
	for _, value := range admins {
		if value.UserID == int(update.Message.From.ID) {
			return true
		}
	}

	return false
}

func CheckAccessForCallback(update tgbotapi.Update, db *sqlx.DB) bool {
	var admins []entities.Admin
	err := db.Select(&admins, "SELECT * FROM admins")
	if err != nil {
		log.Println(err)
	}
	for _, value := range admins {
		if value.UserID == int(update.CallbackQuery.From.ID) {
			return true
		}
	}

	return false
}

func (d *DataBase) NewRowIntoDB(update *tgbotapi.Update) error {
	user := entities.NewUser(update)
	query := fmt.Sprintf("INSERT INTO %s (userid,fname,lname,username,languages,chatid) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (userid) DO NOTHING", usersTable)
	_, err := d.db.Exec(query, user.UserID, user.FName, user.LName, user.UserName, user.Language, user.ChatId)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("INSERT INTO %s (mowner,mdate,inputdata, chatid, chattitle) VALUES ($1,$2,$3,$4,$5)", messagesTable)
	message := entities.NewMessage(update)
	_, err = d.db.Exec(query, user.UserID, message.Date, message.InputData, message.ChatID, message.ChatTitle)
	if err != nil {
		return err
	}
	return nil

}

func (d *DataBase) ShowAllUsers(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {

	if CheckAccess(update, d.db) {

		var users []entities.User
		d.db.Select(&users, "SELECT * FROM users")
		tmpl := template.New("users")
		tmpl, err := tmpl.Parse(allUsersReturn)
		if err != nil {
			return tgbotapi.NewMessage(update.Message.Chat.ID, "Something broke"), err
		}
		var returned buffer.Buffer
		for _, value := range users {
			if err := tmpl.Execute(&returned, value); err != nil {
				return tgbotapi.NewMessage(update.Message.Chat.ID, "Something broke"), err
			}
		}
		result := returned.String()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		return msg, nil

	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You aren't admin of this bot!")
		return msg, nil
	}
}

func (d *DataBase) ShowAllMessages(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {

	if CheckAccess(update, d.db) {
		var messages []entities.Message
		err := d.db.Select(&messages, "SELECT * FROM messages")
		if err != nil {
			return tgbotapi.NewMessage(update.Message.Chat.ID, "Something broke"), err
		}
		tmpl := template.New("messages")
		tmpl, err = tmpl.Parse(allMessagesReturn)
		if err != nil {
			return tgbotapi.NewMessage(update.Message.Chat.ID, "Something broke"), err
		}
		var returnded buffer.Buffer
		for _, value := range messages {
			if err := tmpl.Execute(&returnded, value); err != nil {
				return tgbotapi.NewMessage(update.Message.Chat.ID, "Something broke "), err
			}
		}
		result := returnded.String()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		return msg, nil
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You aren't admin of this bot!")
		return msg, nil
	}
}

func (d *DataBase) ReturnUsers(update tgbotapi.Update) ([]entities.User, error) {

	var users []entities.User

	if CheckAccess(update, d.db) {
		err := d.db.Select(&users, "SELECT * FROM users")
		return users, err
	} else {
		err := errors.New("no admin")
		return users, err
	}
}

func (d *DataBase) ReturnMessages(update tgbotapi.Update) ([]entities.Message, error) {

	var messages []entities.Message

	if CheckAccess(update, d.db) {
		err := d.db.Select(&messages, "SELECT chattitle, chatid FROM messages GROUP BY chattitle, chatid")
		return messages, err
	} else {
		err := errors.New("no admin")
		return messages, err
	}
}

func (d *DataBase) SelectOneUser(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {

	if CheckAccessForCallback(update, d.db) {

		var messages []entities.Message
		request := strings.Split(update.CallbackQuery.Data, "_")
		err := d.db.Select(&messages, "SELECT * FROM messages WHERE mowner = $1", request[1])
		fmt.Println(messages)
		if err != nil {
			return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Something broke "), err
		}
		tmpl := template.New("messages")
		tmpl, err = tmpl.Parse(allMessagesReturn)
		if err != nil {
			return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Something broke"), err
		}
		var returnded buffer.Buffer
		for _, value := range messages {
			if err := tmpl.Execute(&returnded, value); err != nil {
				return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Something broke "), err
			}
		}
		result := returnded.String()
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, result)
		return msg, nil
	} else {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You aren't admin of this bot!")
		return msg, nil
	}
}

func (d *DataBase) SelectOneChat(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {

	if CheckAccessForCallback(update, d.db) {

		var messages []entities.Message
		request := strings.Split(update.CallbackQuery.Data, "_")
		err := d.db.Select(&messages, "SELECT * FROM messages WHERE chatid = $1", request[1])
		if err != nil {
			return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Something broke "), err
		}
		tmpl := template.New("messages2")
		tmpl, err = tmpl.Parse(allMessagesReturn)
		if err != nil {
			return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Something broke"), err
		}
		var returnded buffer.Buffer
		for _, value := range messages {
			if err := tmpl.Execute(&returnded, value); err != nil {
				return tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Something broke "), err
			}
		}
		result := returnded.String()
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, result)
		return msg, nil
	} else {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "You aren't admin of this bot!")
		return msg, nil
	}
}

func (d *DataBase) NewAdmin(update tgbotapi.Update) error {

	if CheckAccess(update, d.db) {
		var madmin []entities.Admin
		err := d.db.Select(&madmin, "SELECT * FROM admins WHERE masteradmin = true")
		fmt.Println(madmin)
		if err != nil {
			return err
		}
		newAdmin := entities.NewAdminFow(update)
		fmt.Println(newAdmin)
		query := fmt.Sprintf("INSERT INTO %s (userid,username,masteradmin) VALUES ($1,$2,$3)", adminsTable)
		_, err = d.db.Exec(query, newAdmin.UserID, newAdmin.UserName, newAdmin.MasterAdmin)
		return err
	} else {
		return errors.New("you aren't admin of this bot")
	}

}
