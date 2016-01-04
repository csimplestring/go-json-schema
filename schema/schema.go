package schema

import (
	"encoding/json"
	"regexp"
)

//
//type Schema struct {
//	ID string
//	Schema string
//	Title string
//	Description string
//}

// firstly, assume the schema passed in is valid, after the schema validation is finished, you can validate the schema
// against it.

type SchemaInterface interface {
	Validate(v interface{}, path string, interactive bool) validationError
}

type Schema map[string]interface{}

func (s Schema) getFloat64Value(key string) (value float64, exist bool) {
	value = 0
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(json.Number).Float64()
		exist = true
	}

	return
}

func (s Schema) getBoolValue(key string) (value bool, exist bool) {
	value = false
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(bool)
		exist = true
	}

	return
}

func (s Schema) getIntValue(key string) (value int, exist bool) {
	value = 0
	exist = false

	if v, ok := s[key]; ok {
		i64, _ := v.(json.Number).Int64()
		value = int(i64)
		exist = true
	}

	return
}

func (s Schema) getStringValue(key string) (value string, exist bool) {
	value = ""
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(string)
		exist = true
	}

	return
}

// validation keywords for any instance

type Type struct {
	IsArray bool
	Value   JsonType
	Values  []JsonType
}

func (s Schema) Type() (t *Type, exist bool) {
	v, exist := s["type"]
	if !exist {
		return
	}

	t = &Type{}
	switch v.(type) {
	case string:
		t.IsArray = false
		t.Value = JsonType(v.(string))
		return
	case []interface{}:
		t.IsArray = true
		for _, vi := range v.([]interface{}) {
			t.Values = append(t.Values, JsonType(vi.(string)))
		}
		return
	}

	return
}

func (s Schema) Enum() (enums []interface{}, exist bool) {
	v, exist := s["enum"]
	if !exist {
		return
	}

	exist = true
	enums = v.([]interface{})
	return
}

func (s Schema) AllOf() (all []Schema, exist bool) {
	v, exist := s["allOf"]
	if !exist {
		return
	}

	// v must be an array of valid schema
	exist = true
	for _, one := range v.([]interface{}) {
		all = append(all, Schema(one.(map[string]interface{})))
	}
	return
}

func (s Schema) AnyOf() (any []Schema, exist bool) {
	v, exist := s["anyOf"]
	if !exist {
		return
	}

	// v must be an array of valid schema
	exist = true
	for _, one := range v.([]interface{}) {
		any = append(any, Schema(one.(map[string]interface{})))
	}
	return
}

func (s Schema) OneOf() (all []Schema, exist bool) {
	v, exist := s["oneOf"]
	if !exist {
		return
	}

	// v must be an array of valid schema
	exist = true
	for _, one := range v.([]interface{}) {
		all = append(all, Schema(one.(map[string]interface{})))
	}
	return
}

func (s Schema) Not() (not Schema, exist bool) {
	v, exist := s["not"]
	if !exist {
		return
	}

	// v must be an object
	exist = true
	not = Schema(v.(map[string]interface{}))
	return
}

// validation keywords for numeric

func (s Schema) MultipleOf() (divided float64, exist bool) {
	return s.getFloat64Value("multipleOf")
}

func (s Schema) Maximum() (max float64, exist bool) {
	return s.getFloat64Value("maximum")
}

func (s Schema) Minimum() (max float64, exist bool) {
	return s.getFloat64Value("minimum")
}

func (s Schema) ExclusiveMinimum() bool {
	if v, ok := s.getBoolValue("exclusiveMinimum"); ok {
		return v
	}
	return false
}

func (s Schema) ExclusiveMaximum() bool {
	if v, ok := s.getBoolValue("exclusiveMaximum"); ok {
		return v
	}
	return false
}

// validation keywords for string

func (s Schema) MaxLength() (maxLen int, exist bool) {
	return s.getIntValue("maxLength")
}

func (s Schema) MinLength() (minLen int, exist bool) {
	return s.getIntValue("minLength")
}

func (s Schema) Pattern() (pattern string, exist bool) {
	return s.getStringValue("pattern")
}

// validation keywords for array

type AdditionalItems struct {
	IsBool bool
	Bool   bool
	Schema Schema
}

func (s Schema) AdditionalItems() (a *AdditionalItems, exist bool) {
	v, exist := s["additionalItems"]
	if !exist {
		return
	}

	a = &AdditionalItems{}
	switch v.(type) {
	case bool:
		exist = true
		a.IsBool = true
		a.Bool = v.(bool)
		return
	case map[string]interface{}:
		exist = true
		a.IsBool = false
		a.Schema = Schema(v.(map[string]interface{}))
		return
	default:
		return
	}
}

type Items struct {
	IsArray     bool
	ItemSchema  Schema
	ItemSchemas []Schema
}

func (s Schema) Items() (items *Items, exist bool) {
	v, exist := s["items"]
	if !exist {
		return
	}

	// item can be an object or an array of objects
	items = &Items{}
	switch v.(type) {
	case map[string]interface{}:
		exist = true
		items.IsArray = false
		items.ItemSchema = Schema(v.(map[string]interface{}))
		return
	case []interface{}:
		exist = true
		items.IsArray = true
		for _, vi := range v.([]interface{}) {
			items.ItemSchemas = append(items.ItemSchemas, Schema(vi.(map[string]interface{})))
		}
		return
	default:
		return
	}
}

func (s Schema) MaxItems() (maxItems int, exist bool) {
	return s.getIntValue("maxItems")
}

func (s Schema) MinItems() (minItems int, exist bool) {
	return s.getIntValue("minItems")
}

func (s Schema) UniqueItems() bool {
	if v, ok := s.getBoolValue("uniqueItems"); ok {
		return v
	}
	return false
}

// validation keywords for object

func (s Schema) MaxProperties() (maxProperties int, exist bool) {
	return s.getIntValue("maxProperties")
}

func (s Schema) MinProperties() (minProperties int, exist bool) {
	return s.getIntValue("minProperties")
}

func (s Schema) Required() (required []string, exist bool) {
	v, exist := s["required"]
	if !exist {
		return
	}

	for _, r := range v.([]interface{}) {
		required = append(required, r.(string))
	}

	return
}

// properties returns a map where the key is the property name and the value
// is a valid json schema for this property.

type Properties map[string]Schema

func (s Schema) Properties() (properties Properties, exist bool) {
	v, exist := s["properties"]
	if !exist {
		return
	}

	properties = make(Properties)
	for key, val := range v.(map[string]interface{}) {
		properties[key] = Schema(val.(map[string]interface{}))
	}

	return
}

// The value of "additionalProperties" MUST be a boolean or an object.
// If it is an object, it MUST also be a valid JSON Schema.

type AdditionalProperties struct {
	Schema    Schema
	BoolValue bool

	IsBool   bool
	IsSchema bool
}

func (s Schema) AdditionalProperties() (additionalProperties *AdditionalProperties, exist bool) {
	v, exist := s["additionalProperties"]
	if !exist {
		return
	}

	additionalProperties = &AdditionalProperties{}
	switch v.(type) {
	case bool:
		additionalProperties.IsBool = true
		additionalProperties.BoolValue = v.(bool)
		return
	case map[string]interface{}:
		additionalProperties.IsSchema = true
		additionalProperties.Schema = Schema(v.(map[string]interface{}))
		return
	default:
		return
	}
}

// The value of "patternProperties" MUST be a boolean or an object.

type PatternProperties struct {
	patternSchemas map[*regexp.Regexp]Schema
}

func (p *PatternProperties) Match(prop string) (isMatched bool, matchedSchema Schema) {
	for r, s := range p.patternSchemas {
		if r.MatchString(prop) {
			isMatched = true
			matchedSchema = s
			return
		}
	}

	return false, nil
}

func (s Schema) PatternProperties() (patternProperties *PatternProperties, exist bool) {
	v, exist := s["patternProperties"]
	if !exist {
		return
	}

	patternProperties = &PatternProperties{
		patternSchemas: make(map[*regexp.Regexp]Schema),
	}

	for key, val := range v.(map[string]interface{}) {
		r := regexp.MustCompile(key)
		patternProperties.patternSchemas[r] = Schema(val.(map[string]interface{}))
	}

	return
}
