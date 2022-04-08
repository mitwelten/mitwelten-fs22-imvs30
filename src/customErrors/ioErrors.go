package customErrors

import "fmt"

// Invalid output filename
type ErrIOWrite struct {
	err error
}

func (*ErrIOWrite) Error() string {
	return fmt.Sprintf("cannot write to file")
}
