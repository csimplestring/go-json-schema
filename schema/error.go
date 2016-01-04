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
	ArrayItem

	ObjectMaxPropertiesError      = ErrorCode("max properties")
	ObjectMinPropertiesError      = ErrorCode("min properties")
	ObjectRequiredPropertiesError = ErrorCode("required properties")
	ObjectUndefinedPropertyError = ErrorCode("undifined property")

	TypeError          = ErrorCode("type")
	TypeNotMatchError  = ErrorCode("not match type")
	TypesNotMatchError = ErrorCode("not match one of types")

	EnumError = ErrorCode("enum")

	AllOfError = ErrorCode("allOf")
	AnyOfError = ErrorCode("anyOf")
	OneOfError = ErrorCode("oneOf")
	NotError   = ErrorCode("not")

	UndefinedTypeError = ErrorCode("undefined type")
)

type ValidationError interface {
	error
	Code() ErrorCode
	Path() string
}

type validationError struct {
	code ErrorCode
	path string
}

func newError(code ErrorCode, path string) *validationError {
	return &validationError{code, path}
}

func (s *validationError) Code() ErrorCode {
	return s.code
}

func (s *validationError) Path() string {
	return s.path
}

func (s *validationError) Error() string {
	return fmt.Sprintf("Error: %s, Path: %s", s.Code(), s.Path())
}
