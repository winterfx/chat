package rest

import "chat/internal/services/user"

type LoginResponse struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

type RegisterResponse struct {
	Uid string `json:"uid"`
}

type ListFriendResponse struct {
	Uid    string          `json:"uid"`
	Name   string          `json:"name"`
	Status user.UserStatus `json:"status"`
}
