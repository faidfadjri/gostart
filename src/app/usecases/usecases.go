package usecases

import (
	"github.com/faidfadjri/gostart/src/app/usecases/task"
	"github.com/faidfadjri/gostart/src/app/usecases/user"
)

var (
	NewUserUsecase = user.NewUserUsecase
	NewTaskUsecase = task.NewTaskUsecase
)
