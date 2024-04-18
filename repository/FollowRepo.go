package repository

import (
	"context"
	"follower-service/model"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func NewUserRepository(logger *log.Logger) (*UserRepository, error) {
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	return &UserRepository{
		driver: driver,
		logger: logger,
	}, nil
}

func (ur *UserRepository) Follow(user, follower *model.User) error {
	ctx := context.Background()
	session := ur.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	savedUser, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (u:User {userId: $userID}), (f:User {userId: $followerID}) "+
					"CREATE (f)-[:FOLLOWS]->(u)",
				map[string]any{"UserId": user.UserId, "followerID": follower.UserId})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		ur.logger.Println("Error inserting user:", err)
		return err
	}
	ur.logger.Println(savedUser.(string))
	return nil
}

func (ur *UserRepository) Unfollow(user, unfollower *model.User) error {
	ctx := context.Background()
	session := ur.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	savedUser, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (f:User {userId: $followerID})-[r:FOLLOWS]->(u:User {userId: $userID}) "+
					"DELETE r",
				map[string]any{"UserId": user.UserId, "UnfollowerID": unfollower.UserId})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		ur.logger.Println("Error inserting user:", err)
		return err
	}
	ur.logger.Println(savedUser.(string))
	return nil
}
