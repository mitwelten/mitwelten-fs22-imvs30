package customErrors

import "fmt"

// Unfulfilled Min Argument Error
type ErrArgParserUnfulfilledMinArguments struct {
	err error
}

func (*ErrArgParserUnfulfilledMinArguments) Error() string {
	return fmt.Sprintf("expected at least '-input' '-output' and '-method' arguments")
}

// Invalid output filename
type ErrArgParserInvalidOutputFilename struct {
	err error
}

func (*ErrArgParserInvalidOutputFilename) Error() string {
	return fmt.Sprintf("when using -output 'file' a valid -output_fielname must be specified")
}

// Invalid output port
type ErrArgParserInvalidOutputPort struct {
	err error
}

func (*ErrArgParserInvalidOutputPort) Error() string {
	return fmt.Sprintf("when using -output 'stream' a valid -output_port must be specified")
}

// Invalid Argument Error
type ErrArgParserInvalidArgument struct {
	Argument string
	err      error
}

func (e *ErrArgParserInvalidArgument) Error() string {
	return fmt.Sprintf("invalid output argument: -output argument '" + e.Argument + "' not valid. Use -output 'stream' or -output 'file'")
}
