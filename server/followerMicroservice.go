package server

import (
	"context"
	"fmt"
	"follower-service/model"
	"follower-service/service"
	"log"
)

type FollowerMicroservice struct {
	UnimplementedFollowerMicroserviceServer
	FollowerService *service.UserService
}

func (s *FollowerMicroservice) FollowUser(ctx context.Context, in *FollowRequest) (*StringMessage, error) {

	user := &model.User{UserId: int(in.UserID)}
	follower := &model.User{UserId: int(in.FollowerID)}

	err := s.FollowerService.Follow(user, follower)

	if err != nil {
		fmt.Println("Error while following:", err)
		message := &StringMessage{Message: "Error while following"}
		return message, err
	}

	message := &StringMessage{Message: "Successfully followed"}
	return message, err
}
func (s *FollowerMicroservice) UnfollowUser(ctx context.Context, in *FollowRequest) (*StringMessage, error) {
	user := &model.User{UserId: int(in.UserID)}
	follower := &model.User{UserId: int(in.FollowerID)}

	err := s.FollowerService.Unfollow(user, follower)

	if err != nil {
		fmt.Println("Error while unfollowing:", err)
		message := &StringMessage{Message: "Error while unfollowing"}
		return message, err
	}

	message := &StringMessage{Message: "Successfully unfollowed"}
	return message, err
}
func (s *FollowerMicroservice) GetFollowers(ctx context.Context, in *FollowerIdRequest) (*FollowerListResponse, error) {
	user := &model.User{UserId: int(in.Id)}

	followers, err := s.FollowerService.GetFollowers(user)

	if err != nil {
		log.Printf("Error fetching followers", err)
		return nil, err
	}

	var response []*FollowerResponse
	for _, follow := range followers {
		response = append(response, &FollowerResponse{
			Id: int64(follow.UserId),
		})
	}
	return &FollowerListResponse{Followers: response}, nil
}
func (s *FollowerMicroservice) GetFollowings(ctx context.Context, in *FollowerIdRequest) (*FollowerListResponse, error) {
	user := &model.User{UserId: int(in.Id)}

	followers, err := s.FollowerService.GetFollowing(user)

	if err != nil {
		log.Printf("Error fetching followings", err)
		return nil, err
	}

	var response []*FollowerResponse
	for _, follow := range followers {
		response = append(response, &FollowerResponse{
			Id: int64(follow.UserId),
		})
	}
	return &FollowerListResponse{Followers: response}, nil
}
func (s *FollowerMicroservice) GetFollowerSuggestions(ctx context.Context, in *FollowerIdRequest) (*FollowerListResponse, error) {
	user := &model.User{UserId: int(in.Id)}

	followers, err := s.FollowerService.GetFollowerSuggestions(user)

	if err != nil {
		log.Printf("Error fetching followers", err)
		return nil, err
	}

	var response []*FollowerResponse
	for _, follow := range followers {
		response = append(response, &FollowerResponse{
			Id: int64(follow.UserId),
		})
	}
	return &FollowerListResponse{Followers: response}, nil
}
