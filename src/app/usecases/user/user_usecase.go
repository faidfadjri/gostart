package user

// UserUsecase handles HTTP requests
type userUsecase struct {
	// TODO
}

func NewUserUsecase() UserUsecase {
	return &userUsecase{}
}

func (t *userUsecase) DoSomething() error {
	return nil
}
