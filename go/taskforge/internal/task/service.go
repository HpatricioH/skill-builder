package task

type Service struct {
	tasks []Task
}

func NewService() *Service {
	return &Service{
		tasks: []Task{},
	}
}
