package run

type Err string

func (err Err) Error() string {
	return string(err)
}
