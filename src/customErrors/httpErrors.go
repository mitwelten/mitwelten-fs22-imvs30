package customErrors

import "fmt"

// Invalid output filename
type ErrHttpOpenOutputSocket struct {
	err error
}

func (*ErrHttpOpenOutputSocket) Error() string {
	return fmt.Sprintf("can't open HTTP output")
}
