package schema
import (
	"regexp"
	"fmt"
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
	o.validateProperties(obj, path)
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
	properties, exist := o.schema.Properties()
	if !exist {
		return
	}

	patternProperties, hasPatternProperties := o.schema.PatternProperties()
	additionalProperties, hasAdditionalProperties := o.schema.AdditionalProperties()

	for key, val := range obj {
		subPath := fmt.Sprintf("%s.%s", path, key)

		if schema, found := properties[key]; found {
			c := NewBaseConstraint(schema)
			c.Validate(val, subPath)
			o.addErrors(c.Errors())
			continue
		}

		if hasPatternProperties {
			if match, matchedSchema := patternProperties.Match(key); match {
				c := NewBaseConstraint(matchedSchema)
				c.Validate(val, subPath)
				o.addErrors(c.Errors())
				continue
			}
		}

		if hasAdditionalProperties && additionalProperties.IsSchema {
			c := NewBaseConstraint(additionalProperties.Schema)
			c.Validate(val, subPath)
			o.addErrors(c.Errors())
			continue
		}

		if hasAdditionalProperties && additionalProperties.IsBool {
			if additionalProperties.BoolValue == false {
				o.addError(newError(ObjectUndefinedPropertyError, subPath))
			}
		}
	}
}

