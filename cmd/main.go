package main

import (
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/db"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/server"
	"os"
)

func main() {

	readEnv()

	storage, err := db.NewStorage(databaseDSN)
	if err != nil {
		panic(err)
	}

	err = storage.Init()
	if err != nil {
		panic(err)
	}

	s := server.NewServer(":8080", storage)
	err = s.Run()
	if err != nil {
		panic(err)
	}

}

var (
	databaseDSN string
)

func readEnv() {
	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		databaseDSN = envDataBaseDSN
	}
}
