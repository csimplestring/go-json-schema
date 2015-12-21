package schema

import "regexp"

type StringConstraint struct {
	schema Schema
	baseConstraint
}

func NewStringConstraint(schema Schema) *StringConstraint {
	return &StringConstraint{
		schema: schema,
	}
}

func (constraint *StringConstraint) Validate(v interface{}, path string) {
	str := v.(string)
	strLen := len(str)

	if maxLen, ok := constraint.schema.MaxLength(); ok {
		if strLen > maxLen {
			constraint.addError(newError(StringMaxLengthError, path))
		}
	}

	if minLen, ok := constraint.schema.MinLength(); ok {
		if strLen < minLen {
			constraint.addError(newError(StringMinLengthError, path))
		}
	}

	if pattern, ok := constraint.schema.Pattern(); ok {
		if !regexp.MustCompile(pattern).MatchString(str) {
			constraint.addError(newError(StringPatternError, path))
		}
	}
}
