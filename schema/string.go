package schema

type StringConstraint struct {
	schema Schema
	baseConstraint
}

func (constraint *StringConstraint) Validate(v interface{}, path string) {
//	str, ok := v.(string)
//	if !ok {
//		constraint.addError(newError(ErrorCodeStringTypeMismatch, path))
//	}

	// maxLength
//	strLen := len(str)
//	if v, ok := constraint.schema["maxLength"]; ok {
//		maxLen, _ := v.(json.Number).Int64()
//
//	}
}
