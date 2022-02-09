package run

type Task struct {
	//for handle value (V)
	H func(v interface{}) error
	//value for handle (H)
	V interface{}
	//for notice handle (H) error
	C chan error
}

type Tasker interface {
	Add(Task) error
	Stop()
}

// type Handler interface {
// 	Handle(v interface{}) error
// }
