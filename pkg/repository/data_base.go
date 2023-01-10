package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type DataBase struct {
	db *sqlx.DB
}

func NewDataBase(db *sqlx.DB) *DataBase {
	return &DataBase{db: db}
}

const (
	usersTable    = "users"
	messagesTable = "messages"
	adminsTable   = "admins"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Printf("Connected to db is successful")
	query := fmt.Sprintf("INSERT INTO %s (userid,username,masteradmin) VALUES ($1,$2,$3)", adminsTable)
	db.Exec(query, viper.GetString("admin_id"), viper.GetString("admin_name"), true)
	return db, nil
}
