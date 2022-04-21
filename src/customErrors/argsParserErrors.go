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
	return fmt.Sprintf("when using -output 'file' a valid -output_filename must be specified")
}

// Invalid output port
type ErrArgParserInvalidOutputPort struct {
	err error
}

func (*ErrArgParserInvalidOutputPort) Error() string {
	return fmt.Sprintf("when using -output 'stream' a valid -output_port must be specified")
}

// ErrArgParserInvalidMode Invalid method parameter
type ErrArgParserInvalidMode struct {
	Argument string
	err      error
}

func (e *ErrArgParserInvalidMode) Error() string {
	return fmt.Sprintf("invalid method argument: -method argument '" + e.Argument + "' not valid.")
}

// ErrArgParserInvalidGridDimension Invalid method: Grid invalid paramters
type ErrArgParserInvalidGridDimension struct {
	err error
}

func (*ErrArgParserInvalidGridDimension) Error() string {
	return fmt.Sprintf("when using -method 'grid', the grid dimension must be specified. Usage '-grid_dimension \"<row> <col>\"'")
}

// Invalid Argument Error
type ErrArgParserInvalidArgument struct {
	Argument string
	err      error
}

func (e *ErrArgParserInvalidArgument) Error() string {
	return fmt.Sprintf("invalid output argument: -output argument '" + e.Argument + "' not valid. Use -output 'stream' or -output 'file'")
}
