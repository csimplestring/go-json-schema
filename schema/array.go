package schema

import (
	"fmt"
	"reflect"
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
	arr := v.([]interface{})

	constraint.validateMaxItems(arr, path)
	constraint.validateMinItems(arr, path)
	constraint.validateUniqueItem(arr, path)

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

func (constraint *ArrayConstraint) validateMaxItems(items []interface{}, path string) {
	if max, exist := constraint.schema.MaxItems(); exist {
		if len(items) > max {
			constraint.addError(newError(ArrayMaxItemError, path))
		}
	}
}

func (constraint *ArrayConstraint) validateMinItems(items []interface{}, path string) {
	if min, exist := constraint.schema.MinItems(); exist {
		if len(items) < min {
			constraint.addError(newError(ArrayMinItemError, path))
		}
	}
}

func (constraint *ArrayConstraint) validateUniqueItem(items []interface{}, path string) {
	length := len(items)
	if length == 0 {
		return
	}

	one := items[0]
	for i := 1; i < length; i++ {
		if !reflect.DeepEqual(one, items[i]) {
			constraint.addError(newError(ArrayUniqueItemError, path+fmt.Sprintf("[%d]", i)))
		}
	}
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
