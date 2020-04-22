package concurrency

// Task is a simple wrapper to a function that takes no arguments and return a (result, error) pair
type Task interface {
	Run() (interface{}, error)
}

func NewTask(fn func() (interface{}, error)) Task {
	return &task{
		fn: fn,
	}
}

func NewNoErrorTask(fn func() interface{}) Task {
	return NewTask(func() (interface{}, error) {
		return fn(), nil
	})
}

func NewErrorTask(err error) Task {
	return NewTask(func() (interface{}, error) {
		return nil, err
	})
}

type task struct {
	fn func() (interface{}, error)
}

func (t *task) Run() (interface{}, error) {
	return t.fn()
}
