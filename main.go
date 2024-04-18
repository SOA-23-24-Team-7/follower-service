package main

import (
	//"FOLLOWER-SERVICE/controller"

	"follower-service/repository"
	"log"
	"os"
)

func main(){
	port := os.Getenv("PORT") // ako budemo uzimali preko ovih varijabli
	if len(port) == 0 {
		port = "8095" // novi port
	}

	// timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[movie-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[movie-store] ", log.LstdFlags)

	// NoSQL: Initialize Movie Repository store
	store, err := repository.NewUserRepository(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	// defer store.CloseDriverConnection(timeoutContext)
	// store.CheckConnection()

}