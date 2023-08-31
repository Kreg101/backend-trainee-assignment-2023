package main

import (
	"fmt"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/db"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/logger"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/server"
	"os"
)

func main() {

	readEnv()

	// set up logger's output
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		f = os.Stderr
	}

	// create new logger
	log := logger.NewLogger(f)

	// create new Storage and check connection
	storage, err := db.NewStorage(databaseDSN, log)
	if err != nil {
		panic(err)
	}

	// initialize storage
	err = storage.Init()
	if err != nil {
		panic(err)
	}

	// create and run server
	s := server.NewServer(serverHost, storage, log)
	err = s.Run()
	if err != nil {
		panic(err)
	}
}

var (
	// databaseDSN for connection to database
	databaseDSN string

	serverHost = ":8080"

	// logFilePath store path for logger file
	logFilePath string
)

// readEnv gets all needed information from environment variables
func readEnv() {

	//for example: host=db user=postgres dbname=postgres sslmode=disable
	databaseDSN = os.Getenv("DATABASE_DSN")

	// for example: qwerty
	databaseDSN += fmt.Sprintf(" password=%s", os.Getenv("DATABASE_PASSWORD"))

	if envServerHost := os.Getenv("SERVER_HOST"); envServerHost != "" {
		serverHost = envServerHost
	}

	logFilePath = os.Getenv("LOG_FILE_PATH")
}
