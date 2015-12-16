package schema

type Constraint interface {
	Errors() []SchemaError
	Validate(v interface{}, path string) error
}
