package schema

import "reflect"

type Constraint interface {
	Errors() []SchemaError
	Validate(v interface{}, path string)
}

//func ConstraintFactory(schema Schema) ([]Constraint, error) {
//	var constraints []Constraint
//	t, types := schema.Type()
//
//	if t != "" {
//		switch t {
//		case "integer", "number":
//			constraints = append(constraints, NewNumericConstraint(schema))
//		case "string":
//			constraints = append(constraints, NewStringConstraint(schema))
//		case "array":
//			constraints = append(constraints, NewArrayConstraint(schema))
//		default:
//			return nil, fmt.Errorf("Create Constraint Error: unsupported json type %s", jsonType)
//		}
//	} else if types != nil {
//		for _, subtype := types {
//			ConstraintFactory()
//		}
//	}
//
//
//}
//
//func createConstraintByType(t string) {
//
//}

type baseConstraint struct {
	schema Schema
	errors []SchemaError
}

func NewBaseConstraint(schema Schema) *baseConstraint {
	return &baseConstraint{
		schema: schema,
	}
}

func (b *baseConstraint) Errors() []SchemaError {
	return b.errors
}

func (b *baseConstraint) addError(e SchemaError) {
	b.errors = append(b.errors, e)
}

func (b *baseConstraint) addErrors(e []SchemaError) {
	b.errors = append(b.errors, e...)
}

func (b *baseConstraint) Validate(v interface{}, path string) {
	b.validateType(v, path)
	b.validateEnum(v, path)
	b.validateAllOf(v, path)

	//	t, err := getJsonType(v)
	//	if err != nil {
	//		b.addError(newError(UndefinedTypeError, path))
	//	}
	//
	//	var c Constraint
	//	switch t {
	//	case JsonInteger, JsonNumber:
	//		c = NewNumericConstraint(b.schema)
	//	case JsonString:
	//		c = NewStringConstraint(b.schema)
	//	case JsonArray:
	//		c = NewArrayConstraint(b.schema)
	//	}
	//
	//	c.Validate(v, path)
	//	b.addErrors(c.Errors())
}

func (b *baseConstraint) validateType(v interface{}, path string) {
	actualType, err := getJsonType(v)
	if err != nil {
		b.addError(newError(TypeError, path))
	}

	expectedType, expectedTypes, exist := b.schema.Type()
	if !exist {
		return
	}

	// single type
	if expectedType != "" {
		if expectedType != actualType {
			b.addError(newError(TypeNotMatchError, path))
		}
		return
	}

	// mixed type
	for _, t := range expectedTypes {
		if t == actualType {
			return
		}
	}
	b.addError(newError(TypesNotMatchError, path))
}

func (b *baseConstraint) validateEnum(v interface{}, path string) {
	enums, exist := b.schema.Enum()
	if !exist {
		return
	}

	for _, enum := range enums {
		if reflect.DeepEqual(enum, v) {
			return
		}
	}
	b.addError(newError(EnumError, path))
}

func (b *baseConstraint) validateAllOf(v interface{}, path string) {
	all, exist := b.schema.AllOf()
	if !exist {
		return
	}

	for _, one := range all {
		c := NewBaseConstraint(one)
		c.Validate(v, path)

		if len(c.Errors()) > 0 {
			b.addError(newError(AllOfError, path))
		}
	}
}

func (b *baseConstraint) validateAnyOf(v interface{}, path string) {

}

func (b *baseConstraint) validateOneOf(v interface{}, path string) {

}

func (b *baseConstraint) validateNot(v interface{}, path string) {

}
