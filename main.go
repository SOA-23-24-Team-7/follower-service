package main

import (
	//"FOLLOWER-SERVICE/controller"

	"follower-service/controller"
	"follower-service/repository"
	"follower-service/service"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	serviceLogger := log.New(os.Stdout, "[movie-service] ", log.LstdFlags)
	service := service.NewUserService(store,serviceLogger)

	//kontroler
	controllerLogger := log.New(os.Stdout, "[movie-controller] ", log.LstdFlags)
	controller := controller.NewUserController(service,controllerLogger)
	if(controller!= nil){

	}

	// endpoints
	router := mux.NewRouter().StrictSlash(true)

	// endpoints for following
	router.HandleFunc("/followers/follow/{userID}/{followerID}", controller.FollowUser).Methods("POST")
	router.HandleFunc("/followers/unfollow/{userID}/{followerID}", controller.UnfollowUser).Methods("POST")
	router.HandleFunc("/followers/getFollowers/{userID}", controller.GetFollowers).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	println("Server starting")
	log.Fatal(http.ListenAndServe(":8095", router)) 
}