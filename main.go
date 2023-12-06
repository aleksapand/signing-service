package main

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	db := persistence.GetInMemoryDB()
	server := api.NewServer(ListenAddress, db)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
