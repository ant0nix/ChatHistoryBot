<b>Описание телеграмм-бота</b>

<b>Цель</b>  
	Я создавал данного бота, чтобы продемонстрировать свою способность писать программы на языке Go, которые могут взаимодействовать с базой данных.  

<b>Описание работы</b>  
	Данный бот сохраняет все текстовые сообщения, которые бот смог обнаружить. Бот может записывать сообщения, которые были отправлены как в лично боту, так и сообщения, которые были отправлены в беседе, где есть данный бот.  
	Также у этого бота есть режим администратора, в котором можно смотреть все сообщения, всех пользователей, которые записал бот.  

<b>Режим администратора</b>  
	Чтобы перейти в режим администратора необходимо написать боту команду /start. При этом вы должны быть администратором. Чтобы стать администратором, напишите мне в телеграмм @Ant0niX. После запуска бота вы уведите все возможные команды.  
  
<b>При создании бота я использовал</b>  
Go — v1.18  
PostgreSQL — v15.1  
  
Библиотеки:    
tgbotapi — v5  
sqlx — v1.3.5  
pq — v1.10.7  
viper — v1.14.0  
gotenv — v1.4.1  
zap — v1.21.0  
