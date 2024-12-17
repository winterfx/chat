package repository

import (
	"chat/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func initDatabase() {
	var err error
	fmt.Printf("connect string: %s\n", config.Config.MySqlConnectString)
	database, err = sql.Open("mysql", config.Config.MySqlConnectString) //todo read from env
	if err != nil {
		panic(err)
	}
}

func GetDefaultUserRepo() UserRepo {
	if database == nil {
		initDatabase()
	}
	return NewUserRepo(database)
}

func GetDefaultChatHistoryRepo() ChatHistoryRepo {
	if database == nil {
		initDatabase()
	}
	return NewChatHistoryRepo(database)
}
