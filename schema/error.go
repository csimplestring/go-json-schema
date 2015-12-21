package schema

import "fmt"

type ErrorCode string

const (
	NumericMultipleOfError       = ErrorCode("multipleOf")
	NumericTypeMismatchError     = ErrorCode("number type")
	NumericMaximumError          = ErrorCode("maximum")
	NumericExclusiveMaximumError = ErrorCode("exclusiveMaximum")
	NumericMinimumError          = ErrorCode("minimum")
	NumericExclusiveMinimumError = ErrorCode("exclusiveMinimum")

	StringTypeMismatchError = ErrorCode("string type")
	StringMinLengthError    = ErrorCode("minLength")
	StringMaxLengthError    = ErrorCode("maxLength")
	StringPatternError      = ErrorCode("pattern")

	ArrayTypeMismatchError   = ErrorCode("array type")
	ArrayMaxItemError        = ErrorCode("maxItems")
	ArrayMinItemError        = ErrorCode("minItems")
	ArrayUniqueItemError     = ErrorCode("uniqueItem")
	ArrayAdditionalItemError = ErrorCode("additionalItem")
	ArrayItemError           = ErrorCode("item")

	TypeError          = ErrorCode("type")
	TypeNotMatchError  = ErrorCode("not match type")
	TypesNotMatchError = ErrorCode("not match one of types")

	EnumError = ErrorCode("enum")

	AllOfError = ErrorCode("allOf")

	UndefinedTypeError = ErrorCode("undefined type")
)

type SchemaError interface {
	error
	Code() ErrorCode
	Path() string
}

type schemaError struct {
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
	return fmt.Sprintf("Error: %s, Path: %s", s.Code(), s.Path())
}
