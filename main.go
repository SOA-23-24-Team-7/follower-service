package main

import (
	"encoding/json"
	"follower-service/model"
	"follower-service/repository"
	"follower-service/server"
	"follower-service/service"
	"log"
	"net"
	"os"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8095"
	}
	logger := log.New(os.Stdout, "[follower-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[follower-store] ", log.LstdFlags)

	store, err := repository.NewUserRepository(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}

	serviceLogger := log.New(os.Stdout, "[follower-service] ", log.LstdFlags)
	service := service.NewUserService(store, serviceLogger)

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	server.RegisterFollowerMicroserviceServer(grpcServer, &server.FollowerMicroservice{
		FollowerService: service,
	})

	conn := Conn()
	defer conn.Close()
	handleNATSSubscription(conn, service)

	listener, err := net.Listen("tcp", ":8095")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server listening on port :8095")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
func handleNATSSubscription(natsConn *nats.Conn, followerService *service.UserService) {
	natsConn.Subscribe("comment.created", func(m *nats.Msg) {
		var event struct {
			UserID    int `json:"user_id"`
			AuthorID  int `json:"author_id"`
			CommentID int `json:"comment_id"`
		}
		log.Printf("SUBSCRIBING NATS REQ  %d,%d", event.AuthorID, event.UserID)
		json.Unmarshal(m.Data, &event)
		author := &model.User{UserId: event.AuthorID}
		user := &model.User{UserId: event.UserID}
		// Attempt to create follow action
		err := followerService.Follow(user, author)
		if err != nil {
			log.Printf("Failed to create follow: %v", err)

			// Send rollback message if follow creation fails
			rollbackEvent := struct {
				CommentID int `json:"comment_id"`
			}{
				CommentID: event.CommentID,
			}
			rollbackData, _ := json.Marshal(rollbackEvent)
			log.Printf("PUBLISH FOLLOWER SERVICE")
			natsConn.Publish("comment.creation.rollback", rollbackData)
		} else {
			log.Printf("Successfully created follow for UserID: %d by AuthorID: %d", event.UserID, event.AuthorID)
		}
	})
}
