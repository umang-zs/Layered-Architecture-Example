package errors

type TestError struct {
	Err error
}

func (m TestError) Error() string {
	return m.Err.Error()
}
