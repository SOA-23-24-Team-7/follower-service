package main

import (
	//"FOLLOWER-SERVICE/controller"

	//"follower-service/controller"

	"follower-service/repository"
	"follower-service/server"
	"follower-service/service"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func main() {
	port := os.Getenv("PORT") // ako budemo uzimali preko ovih varijabli
	if len(port) == 0 {
		port = "8095" // novi port
	}

	//timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[follower-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[follower-store] ", log.LstdFlags)

	//NoSQL: Initialize follower Repository store
	store, err := repository.NewUserRepository(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	//defer store.CloseDriverConnection(timeoutContext)
	//store.CheckConnection()

	serviceLogger := log.New(os.Stdout, "[follower-service] ", log.LstdFlags)
	service := service.NewUserService(store, serviceLogger)

	//kontroler
	//controllerLogger := log.New(os.Stdout, "[follower-controller] ", log.LstdFlags)
	//controller := controller.NewUserController(service, controllerLogger)

	// endpoints
	//router := mux.NewRouter().StrictSlash(true)

	// endpoints for following
	// router.HandleFunc("/followers/follow/{userID}/{followerID}", controller.FollowUser).Methods("POST")
	// router.HandleFunc("/followers/unfollow/{userID}/{followerID}", controller.UnfollowUser).Methods("POST")
	// router.HandleFunc("/followers/getFollowers/{userID}", controller.GetFollowers).Methods("GET")
	// router.HandleFunc("/followers/getFollowing/{userID}", controller.GetFollowings).Methods("GET")
	// router.HandleFunc("/followers/suggestions/{userID}", controller.GetFollowerSuggestions).Methods("GET")

	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	// println("Server starting")
	// log.Fatal(http.ListenAndServe(":8095", router))

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	server.RegisterFollowerMicroserviceServer(grpcServer, &server.FollowerMicroservice{
		FollowerService: service,
	})

	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server listening on port :8088")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
