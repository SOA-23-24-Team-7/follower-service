package service

import (
	"follower-service/model"
	"follower-service/repository"
	"log"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *log.Logger
}

func NewUserService(repo *repository.UserRepository, logger *log.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (us *UserService) Follow(user, follower *model.User) error {
	err := us.repo.Follow(user, follower)
	if err != nil {
		us.logger.Println("Error following user:", err)
		return err
	}
	us.logger.Println("User followed successfully")
	return nil
}

func (us *UserService) Unfollow(user, unfollower *model.User) error {
	err := us.repo.Unfollow(user, unfollower)
	if err != nil {
		us.logger.Println("Error unfollowing user:", err)
		return err
	}
	us.logger.Println("User unfollowed successfully")
	return nil
}

func (us *UserService) GetFollowers(user *model.User) ([]*model.User, error) {
	followers, err := us.repo.GetFollowers(user)
	if err != nil {
		us.logger.Println("Error getting followers:", err)
		return nil, err
	}
	us.logger.Println("Retrieved followers successfully")
	return followers, nil
}

func (us *UserService) GetFollowing(user *model.User) ([]*model.User, error) {
	followers, err := us.repo.GetFollowing(user)
	if err != nil {
		us.logger.Println("Error getting followers:", err)
		return nil, err
	}
	us.logger.Println("Retrieved followers successfully")
	return followers, nil
}
