package repository

import (
	"context"
	"fmt"
	"follower-service/model"
	"log"

	//"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository struct {
	driver       neo4j.DriverWithContext
	logger       *log.Logger
	databaseName string
}

func NewUserRepository(logger *log.Logger) (*UserRepository, error) {
	// uri := os.Getenv("NEO4J_DB")
	// user := os.Getenv("NEO4J_USERNAME")
	// pass := os.Getenv("NEO4J_PASS")
	// auth := neo4j.BasicAuth(user, pass, "")
	uri := "bolt://localhost:7687"
	user := "neo4j"
	pass := "password"
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	return &UserRepository{
		driver:       driver,
		logger:       logger,
		databaseName: "neo4j",
	}, nil
}

func (ur *UserRepository) Follow(user, follower *model.User) error {
	ctx := context.Background()
	session := ur.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ur.databaseName})
	defer session.Close(ctx)
	println(user.UserId)
	savedUser, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MERGE (u:User {userId: $userId}) "+
				"MERGE (f:User {userId: $followerId}) "+
				"CREATE (f)-[:FOLLOWS]->(u)",
				map[string]any{"userId": user.UserId, "followerId": follower.UserId})
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
	ur.logger.Println(savedUser)
	return nil
}

func (ur *UserRepository) Unfollow(user, unfollower *model.User) error {
	ctx := context.Background()
	session := ur.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ur.databaseName})
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

func (ur *UserRepository) GetFollowers(user *model.User) ([]*model.User, error) {
	ctx := context.Background()
	session := ur.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ur.databaseName})
	defer session.Close(ctx)

	f, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (f:User)-[:FOLLOWS]->(u:User {userId: $userID}) RETURN f",
				map[string]any{"userID": user.UserId})
			if err != nil {
				println(err)
				return nil, err
			}

			var followers []*model.User
			for result.Next(ctx) {
				record := result.Record()
				follower, ok := record.Get("f")
				
				if !ok {
					continue
				}
				followers = append(followers, follower.(*model.User)) 
			}
			return followers, result.Err()

		})

	if err != nil {
		ur.logger.Println("Error querying search:", err)
		return nil, err
	}

	followers, ok := f.([]*model.User)
	if !ok {
		return nil, fmt.Errorf("failed to type assert followers to []*model.User")
	}

	return followers, nil
}


func (ur *UserRepository) GetFollowing(user *model.User) ([]*model.User, error) {
	ctx := context.Background()
	session := ur.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ur.databaseName})
	defer session.Close(ctx)

	f, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (u:User {userId: $userID})-[:FOLLOWS]->(f:User) RETURN f",
				map[string]any{"User": user})
			if err != nil {
				return nil, err
			}

			var followings []*model.User
			for result.Next(ctx) {
				record := result.Record()
				follower, ok := record.Get("f")
				if !ok {
					continue
				}
				followings = append(followings, follower.(*model.User))
			}
			return followings, result.Err()

		})

	if err != nil {
		ur.logger.Println("Error querying search:", err)
		return nil, err
	}
	return f.([]*model.User), nil

}
