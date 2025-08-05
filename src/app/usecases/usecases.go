package usecases

import (
	"github.com/faidfadjri/gostart/src/app/usecases/task"
	"github.com/faidfadjri/gostart/src/app/usecases/user"
)

type (
	TaskUsecase = task.TaskUsecase
	UserUsecase = user.UserUsecase
)

var (
	NewTaskUsecase = task.NewTaskUsecase
	NewUserUsecase = user.NewUserUsecase
)
