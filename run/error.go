package run

type Err string

func (err Err) Error() string {
	return string(err)
}

func FrontErr(errFuncs ...func() error) (err error) {
	for _, f := range errFuncs {
		if err = f(); err != nil {
			break
		}
	}
	return
}

// func MixErr(errs ...error) (err error) {

// }
