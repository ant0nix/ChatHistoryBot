package repository

var allUsersReturn = `
ID: {{.ID}}
UserID: {{.UserID}}
Name: {{.FName}} {{.LName}}
Username: {{.UserName}}
Language: {{.Language}}
`

var allMessagesReturn = `
ID: {{.ID}}
Chat Title: {{.ChatTitle}}
Owner: {{.Owner}}
Date: {{.Date}}
Message: {{.InputData}}
`
var StartMessage = `
Приветсвутю! Это учебный бот-проект, который сохраняет в базу данных весь текст и информацию о пользователях в чате.
Бот одинаково работает как в группе, так и при личной переписке с ним.

<b>ВНИМАНИЕ! Команды, которые показывают базу данных не работают, если вас нет в списке админов. Напишите @Ant0niX, чтобы получить доступ к базе данных!</b>

Команды:
/start - начало работы бота (выводится это сообщение)
/show_all_users - выводит всех пользователей, которых бот обнаружил в беседах, или которые писали боту
/show_all_messages - выводит все сообщения, которые бот получал
/show_users_messages - выводит сообщения одного юзера`
