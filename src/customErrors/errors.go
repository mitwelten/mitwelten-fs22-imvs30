package customErrors

import "fmt"

// Unfulfilled Min Argument Error
type ArgParserUnfulfilledMinArgumentsError struct {
	err error
}

func (*ArgParserUnfulfilledMinArgumentsError) Error() string {
	return fmt.Sprintf("expected at least '-input' '-output' and '-method' arguments")
}

//
type ArgParserInvalidInputError struct {
	err error
}

func (*ArgParserInvalidInputError) Error() string {
	return fmt.Sprintf("-output 'stream' only valid in combination with -output_port ")
}

type ArgParserInvalidOutputActionError struct {
	err error
}

type ArgParserInvalidOutputFilenameError struct {
	err error
}

type ArgParserInvalidOutputPortError struct {
	err error
}

type ArgParserInvalidArguments struct {
	err error
}
