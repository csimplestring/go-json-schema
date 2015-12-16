package schema
import (
	"encoding/json"
	"math"
	"fmt"
	"reflect"
)

type NumericConstraint struct {
	schema Schema
	errors []SchemaError
}

func NewNumericConstraint(schema Schema) *NumericConstraint {
	return &NumericConstraint{
		schema: schema,
	}
}

func (constraint *NumericConstraint) Errors() []SchemaError  {
	return constraint.errors
}

func (constraint *NumericConstraint) addError(e SchemaError)  {
	constraint.errors = append(constraint.errors, e)
}

func (constraint *NumericConstraint) Validate(v interface{}, path string) error {
	schema := constraint.schema

	n, ok := v.(json.Number)
	if !ok {
		return fmt.Errorf("require a json number, but %s provided", reflect.TypeOf(v).String())
	}

	f, err := n.Float64()
	if err != nil {
		return err
	}

	if m, ok := schema["multipleOf"]; ok {
		divided, _ := m.(json.Number).Float64()
		if math.Mod(f, divided) != float64(0) {
			constraint.addError(newError(ErrorCodeMultipleOf, path))
		}
	}

	if v, ok := schema["maximum"]; ok {
		maximum, _ := v.(json.Number).Float64()
		if f > maximum {
			constraint.addError(newError(ErrorCodeMaximum, path))
		}

		if v, ok := schema["exclusiveMaximum"]; ok {
			if v.(bool) && f == maximum {
				constraint.addError(newError(ErrorCodeExclusiveMaximum, path))
			}
		}
	}

	if v, ok := schema["minimum"]; ok {
		minimum , _ := v.(json.Number).Float64()
		if f < minimum {
			constraint.addError(newError(ErrorCodeMinimum, path))
		}

		if v, ok := schema["exclusiveMinimum"]; ok {
			if v.(bool) && f == minimum {
				constraint.addError(newError(ErrorCodeExclusiveMinimum, path))
			}
		}
	}

	return nil
}