package main

import (
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/db"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/server"
)

func main() {

	storage, err := db.NewStorage("host=localhost user=postgres password=Kravchenko01 dbname=really sslmode=disable")
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
