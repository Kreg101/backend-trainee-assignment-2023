package main

import (
	"fmt"
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

	s := server.NewServer(serverHost, storage)
	err = s.Run()
	if err != nil {
		panic(err)
	}

}

var (
	databaseDSN string
	serverHost  string
)

func readEnv() {
	if envDataBaseDSN := os.Getenv("DATABASE_DSN"); envDataBaseDSN != "" {
		databaseDSN = envDataBaseDSN
	}
	databaseDSN += fmt.Sprintf(" password=%s", os.Getenv("DATABASE_PASSWORD"))
	if envServerHost := os.Getenv("SERVER_HOST"); envServerHost != "" {
		serverHost = envServerHost
	}
}
