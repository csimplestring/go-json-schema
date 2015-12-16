package schema

type Constraint interface {
	Errors() []SchemaError
	Validate(v interface{}, path string) error
}

type baseConstraint struct {
	errors []SchemaError
}

func (b *baseConstraint) Errors() []SchemaError {
	return b.errors
}

func (b *baseConstraint) addError(e SchemaError)  {
	b.errors = append(b.errors, e)
}

