package schema
import (
	"fmt"
	"regexp"
)

type ObjectConstraint struct {
	schema       Schema
	propPatterns []*regexp.Regexp
	baseConstraint
}

func NewObjectConstraint(s Schema) *ObjectConstraint {
	return &ObjectConstraint{
		schema: s,
	}
}

func (o *ObjectConstraint) Validate(v interface{}, path string) {
	obj := v.(map[string]interface{})

	o.validateMaxProperties(obj, path)
	o.validateMinProperties(obj, path)
	o.validateRequired(obj, path)
}

func (o *ObjectConstraint) validateMaxProperties(obj map[string]interface{}, path string) {
	max, exist := o.schema.MaxProperties()
	if exist && len(obj) > max {
		o.addError(newError(ObjectMaxPropertiesError, path))
	}
}

func (o *ObjectConstraint) validateMinProperties(obj map[string]interface{}, path string) {
	min, exist := o.schema.MinProperties()
	if exist && len(obj) < min {
		o.addError(newError(ObjectMinPropertiesError, path))
	}
}

func (o *ObjectConstraint) validateRequired(obj map[string]interface{}, path string) {
	required, exist := o.schema.Required()
	if exist {
		for _, prop := range required {
			if _, ok := obj[prop]; !ok {
				o.addError(newError(ObjectRequiredPropertiesError, path))
			}
		}
	}
}

func (o *ObjectConstraint) validateProperties(obj map[string]interface{}, path string) {
	propSchema, exist := o.schema.Properties()
	if !exist {
		return
	}

	patternProperties, hasPatternProperties := o.schema.PatternProperties()
	additionSchema, allowAddition, hasAddition := o.schema.AdditionalProperties()

	for prop, val := range obj {
		subPath := fmt.Sprintf("%s.%s", path, prop)

		var s Schema

		// defined property
		if s, ok := propSchema[prop]; ok {
			c := NewBaseConstraint(s)
			c.Validate(val, subPath)
			o.addErrors(c.Errors())
			continue
		}

		// try pattern property
		if hasPatternProperties && patternProperties.match(prop, val) {
			continue
		}

		// no additional
		if !hasAddition || (hasAddition && !allowAddition ){
			o.addError(newError(ObjectUndefinedPropertyError, subPath))
			continue
		}

		if hasAddition && allowAddition {
			continue
		}

		if hasAddition &&
	}
}

func (o *ObjectConstraint) tryMatchPattern(p string) bool {
	if len(o.propPatterns) == 0 {
		return false
	}

	for _, pp := range o.propPatterns {
		if pp.MatchString(p) {
			return true
		}
	}

	return false
}


func (o *ObjectConstraint) validatePatternProperties(obj map[string]interface{}, path string) {

}

func (o *ObjectConstraint) validateAdditionalProperties(obj map[string]interface{}, path string) {

}

