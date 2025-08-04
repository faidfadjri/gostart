package blog

// BlogUsecase handles HTTP requests
type blogUsecase struct {
	// TODO
}

func NewBlogUsecase() BlogUsecase {
	return &blogUsecase{}
}

func (t *blogUsecase) DoSomething() error {
	return nil
}
