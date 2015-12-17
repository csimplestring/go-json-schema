package schema
import "fmt"

type ErrorCode int

const (

	NumericMultipleOfError ErrorCode = iota
	NumericTypeMismatchError
	NumericMaximumError
	NumericExclusiveMaximumError
	NumericMinimumError
	NumericExclusiveMinimumError

	StringTypeMismatchError
	StringMinLengthError
	StringMaxLengthError
	StringPatternError
)

type SchemaError interface {
	error
	Code() ErrorCode
	Path() string
}

type schemaError struct  {
	code ErrorCode
	path string
}

func newError(code ErrorCode, path string) *schemaError {
	return &schemaError{code, path}
}

func (s *schemaError) Code() ErrorCode {
	return s.code
}

func (s *schemaError) Path() string {
	return s.path
}

func (s *schemaError) Error() string {
	return fmt.Sprintf("%d %s", s.Code, s.Path)
}

