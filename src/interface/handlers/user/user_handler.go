package handler

import "github.com/faidfadjri/gostart/src/app/usecases"

// UserHandler handles HTTP requests
type UserHandler struct {
	usecase usecases.UserUsecase
}

func NewUserHandler(u usecases.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: u,
	}
}

// TODO
