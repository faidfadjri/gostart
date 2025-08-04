package usecases

import (
	"github.com/faidfadjri/gostart/src/app/usecases/blog"
	"github.com/faidfadjri/gostart/src/app/usecases/task"
	"github.com/faidfadjri/gostart/src/app/usecases/user"
)

var (
	NewBlogUsecase = blog.NewBlogUsecase

	NewTaskUsecase = task.NewTaskUsecase

	NewUserUsecase = user.NewUserUsecase
)
