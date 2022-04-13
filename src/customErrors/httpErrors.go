package customErrors

import "fmt"

// Invalid output filename
type ErrHttpOpenOutputSocket struct {
	err error
}

func (*ErrHttpOpenOutputSocket) Error() string {
	return fmt.Sprintf("can't open HTTP output")
}

type ErrHttpOpenInputSocketDial struct {
	IP  string
	err error
}

func (e *ErrHttpOpenInputSocketDial) Error() string {
	return fmt.Sprintf("dial input socket failed for IP address %s", e.IP)
}

type ErrHttpWriteHeader struct {
	IP  string
	err error
}

func (e *ErrHttpWriteHeader) Error() string {
	return fmt.Sprintf("error while writing header to IP address %s", e.IP)
}

type ErrHttpReadHeader struct {
	err error
}

func (*ErrHttpReadHeader) Error() string {
	return fmt.Sprintf("error while reading header")
}

type ErrHttpEmptyFrame struct {
	err error
}

func (*ErrHttpEmptyFrame) Error() string {
	return fmt.Sprintf("error received empty frame")
}

type ErrHttpReadFrame struct {
	err error
}

func (*ErrHttpReadFrame) Error() string {
	return fmt.Sprintf("error while reading full frame")
}

type ErrHttpReadEntireFrame struct {
	err error
}

func (*ErrHttpReadEntireFrame) Error() string {
	return fmt.Sprintf("error while reading frame: Cannot read all bytes")
}
