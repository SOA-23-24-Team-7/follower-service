package controller

import (
	"encoding/json"
	"follower-service/model"
	"follower-service/service"
	"log"
	"net/http"
	"strconv"

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
	userIDStr := vars["userID"]

	followerIDStr := vars["followerID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		uc.logger.Println("Error converting userID to int:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(userID)
	followerID, err := strconv.Atoi(followerIDStr)
	if err != nil {
		uc.logger.Println("Error converting followerID to int:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &model.User{UserId: userID}
	follower := &model.User{UserId: followerID}

	err = uc.userService.Follow(user, follower)
	if err != nil {
		uc.logger.Println("Error following user:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (uc *UserController) UnfollowUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	followerIDStr := vars["followerID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		uc.logger.Println("Error converting userID to int:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	unfollowerID, err := strconv.Atoi(followerIDStr)
	if err != nil {
		uc.logger.Println("Error converting followerID to int:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &model.User{UserId: userID}
	unfollower := &model.User{UserId: unfollowerID}

	err = uc.userService.Unfollow(user, unfollower)
	if err != nil {
		uc.logger.Println("Error unfollowing user:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (uc *UserController) GetFollowers(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		uc.logger.Println("Error converting userID to int:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &model.User{UserId: userID}

	followers, err := uc.userService.GetFollowers(user)

	if err != nil {
		uc.logger.Println("Error getting followers:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if followers == nil {
		followers = []*model.User{}
	}

	json.NewEncoder(rw).Encode(followers)
}

func (uc *UserController) GetFollowings(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		uc.logger.Println("Error converting userID to int:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &model.User{UserId: userID}

	followers, err := uc.userService.GetFollowing(user)
	if err != nil {
		uc.logger.Println("Error getting followings:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if followers == nil {
		followers = []*model.User{}
	}

	json.NewEncoder(rw).Encode(followers)
}

func (uc *UserController) GetFollowerSuggestions(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		uc.logger.Println("Error converting userID to int:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &model.User{UserId: userID}

	followers, err := uc.userService.GetFollowerSuggestions(user)
	if err != nil {
		uc.logger.Println("Error getting follower suggestions:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if followers == nil {
		followers = []*model.User{}
	}

	json.NewEncoder(rw).Encode(followers)
}
