package controller

import (
	"encoding/json"
	"follwer-service/model"
	"follwer-service/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserController struct {
	userService *service.UserService
	logger      *log.Logger
}

func NewUserController(userService *service.UserService, logger *log.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

func (uc *UserController) FollowUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	followerID := vars["followerID"]

	user := &model.User{UserID: userID}
	follower := &model.User{UserID: followerID}

	err := uc.userService.Follow(user, follower)
	if err != nil {
		uc.logger.Println("Error following user:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (uc *UserController) UnfollowUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	unfollowerID := vars["unfollowerID"]

	user := &model.User{UserID: userID}
	unfollower := &model.User{UserID: unfollowerID}

	err := uc.userService.Unfollow(user, unfollower)
	if err != nil {
		uc.logger.Println("Error unfollowing user:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (uc *UserController) GetFollowers(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	user := &model.User{UserID: userID}

	followers, err := uc.userService.GetFollowers(user)
	if err != nil {
		uc.logger.Println("Error getting followers:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(followers)
}
