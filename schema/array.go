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
	constraint.validateItems(arr, path)
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

func (constraint *ArrayConstraint) validateItems(values []interface{}, path string) {
	items, exist := constraint.schema.Items()
	if !exist {
		return
	}

	// list validation
	if !items.IsArray {
		c := NewBaseConstraint(items.ItemSchema)
		for i, v := range values {
			c.Validate(v, fmt.Sprintf("%s[%d]", path, i))
		}
		constraint.addErrors(c.Errors())
		return
	}

	// tuple validation
	additionItems, existAddition := constraint.schema.AdditionalItems()
	itemSchemas := items.ItemSchemas
	itemSchemaSize := len(itemSchemas)

	for i, v := range values {
		subPath := fmt.Sprintf("%s[%d]", path, i)

		if i >= itemSchemaSize {
			// additional schema not exists
			if !existAddition {
				constraint.addError(newError(ArrayAdditionalItemError, subPath))
				continue
			}

			// additional schema is object
			if existAddition && !additionItems.IsBool {
				c := NewBaseConstraint(additionItems.Schema)
				c.Validate(v, subPath)
				constraint.addErrors(c.Errors())
				continue
			}

			// additional schema is true
			if existAddition && additionItems.IsBool && additionItems.Bool == true {
				continue
			}

			// additional schema is false
			if existAddition && additionItems.IsBool && additionItems.Bool == false {
				constraint.addError(newError(ArrayAdditionalItemError, subPath))
				continue
			}
		}

		c := NewBaseConstraint(itemSchemas[i])
		c.Validate(v, subPath)
		constraint.addErrors(c.Errors())
	}
}
