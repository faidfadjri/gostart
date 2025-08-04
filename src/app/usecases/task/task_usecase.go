package task

// TaskUsecase handles HTTP requests
type taskUsecase struct {
	// TODO
}

func NewTaskUsecase() TaskUsecase {
	return &taskUsecase{}
}

func (t *taskUsecase) DoSomething() error {
	return nil
}
