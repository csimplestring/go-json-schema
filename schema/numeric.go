package schema
import (
	"encoding/json"
	"math"
)

type NumericConstraint struct {
	schema Schema
	baseConstraint
}

func NewNumericConstraint(schema Schema) *NumericConstraint {
	return &NumericConstraint{
		schema: schema,
	}
}

func (constraint *NumericConstraint) Validate(v interface{}, path string) {
	schema := constraint.schema

	n, ok := v.(json.Number)
	if !ok {
		constraint.addError(newError(NumericTypeMismatchError, path))
	}

	f, _ := n.Float64()

	if divided, ok := schema.MultipleOf(); ok {
		if math.Mod(f, divided) != float64(0) {
			constraint.addError(newError(NumericMultipleOfError, path))
		}
	}

	if max, ok := schema.Maximum(); ok {
		if f > max {
			constraint.addError(newError(NumericMaximumError, path))
		}

		if schema.ExclusiveMaximum() && f == max {
			constraint.addError(newError(NumericExclusiveMaximumError, path))
		}
	}

	if min, ok := schema.Minimum(); ok {
		if f < min {
			constraint.addError(newError(NumericMinimumError, path))
		}

		if schema.ExclusiveMinimum() && f == min {
			constraint.addError(newError(NumericExclusiveMinimumError, path))
		}
	}
}