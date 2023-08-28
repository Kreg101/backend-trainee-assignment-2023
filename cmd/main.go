package main

import "github.com/Kreg101/backend-trainee-assignment-2023/internal/server"

func main() {

	s := server.NewServer(":8080")
	s.Start()
}
