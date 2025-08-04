package task

// TaskUsecase handles HTTP requests
type TaskUsecase struct {
	usecase *TaskUsecase
}

func NewTaskUsecase(u *TaskUsecase) *TaskUsecase {
	return &TaskUsecase{
		usecase: u,
	}
}
