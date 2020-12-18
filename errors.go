package igbinary

import "fmt"

type DecodeError struct {
	err error
}

func (d DecodeError) Error() string {
	return `igbinary: Decode(` + d.err.Error() + `)`
}

func decodeError(err error) error {
	return &DecodeError{
		err: err,
	}
}

func decodeErrorF(format string, args ...interface{}) error {
	return &DecodeError{
		fmt.Errorf(format, args...),
	}
}
