package schema
import (
	"reflect"
	"fmt"
)

type ArrayConstraint struct {
	schema Schema
	baseConstraint
}

func NewArrayConstraint(schema Schema) *ArrayConstraint {
	return &ArrayConstraint{
		schema: schema,
	}
}

func (constraint *ArrayConstraint) Validate(v interface{}, path string) {
	arr, ok := v.([]interface{})
	if !ok {
		constraint.addError(newError(ArrayTypeMismatchError, path))
	}

	if !constraint.validateMaxItems(arr, path) {
		return
	}

	if !constraint.validateMinItems(arr, path) {
		return
	}

	if !constraint.validateUniqueItem(arr, path) {
		return
	}

	schema, schemaArray, exist := constraint.schema.Items()
	if !exist {
		return
	}

	if schema != nil {
		// list validation
		return
	}

	if schemaArray != nil {
		additionalItemSchema, boolValue, hasAdditional := constraint.schema.AdditionalItems()
		if !hasAdditional {
			// tuple validation
		}

		if additionalItemSchema != nil {
			// validate additional element
		} else if true == boolValue {
			// always succeeds
		} else if false == boolValue {
			if len(arr) > len(schemaArray) {
				return
			}
		}
	}
}

func (constraint *ArrayConstraint) validateMaxItems(items []interface{}, path string) bool {
	if max, exist := constraint.schema.MaxItems(); exist {
		if len(items) > max {
			constraint.addError(newError(ArrayMaxItemError, path))
			return false
		}
	}

	return true
}

func (constraint *ArrayConstraint) validateMinItems(items []interface{}, path string) bool {
	if min, exist := constraint.schema.MinItems(); exist {
		if len(items) < min {
			constraint.addError(newError(ArrayMinItemError, path))
			return false
		}
	}

	return true
}

func (constraint *ArrayConstraint) validateUniqueItem(items []interface{}, path string) bool {
	length := len(items)
	if length == 0 {
		return true
	}

	one := items[0]
	for i := 1; i < length; i++ {
		if !reflect.DeepEqual(one, items[i]) {
			constraint.addError(newError(ArrayUniqueItemError, path + fmt.Sprintf("[%d]", i)))
			return false
		}
	}

	return true
}
//
//func (constraint *ArrayConstraint) validateItems(items []interface{}, path string) bool {
//	listSchema, tupleSchema, exist := constraint.schema.Items()
//	if !exist {
//		return true
//	}
//}
//
//func (constraint *ArrayConstraint) validateListItems(itemSchema Schema, items []interface{}, path string) bool {
//	for _, item :=range items {
//		constraint.baseConstraint.Validate(itemSchema, item, path)
//	}
//}
