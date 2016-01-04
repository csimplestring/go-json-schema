package constraint
import (
	"reflect"
	"fmt"
	"encoding/json"
)



type Constraint interface {
	Keys() []string
	// look symfony's validation context
	Validate(v interface{}, ctx ValidationContext) []error
}


type schema struct {
	reserved   map[string]Keyword
	customized map[string]Keyword
}

func NewSchema(raw map[string]interface{}) Schema {

}

func (s *schema) Get(name string) (Keyword, exist bool) {
	if k, ok := s.reserved[name]; ok {
		return k, true
	} else if k, ok := s.customized[name]; ok {
		return k, true
	} else {
		return nil, false
	}
}

func (s *schema) Register(keyword Keyword) error {
	name := keyword.Name()

	if _, exist := s.customized[name]; exist {
		return fmt.Errorf("Keyword %s already exists", name)
	}

	s.customized[name] = keyword
	return nil
}

type SchemaConstraint struct {
	schema Schema
}

func (s *SchemaConstraint) Keys() []string {
	return []string{"type"}
}

func (s *SchemaConstraint) Validate(v interface{}, ctx ValidationContext) (errors []error) {
	for _, key := range s.Keys() {
		if k, exist := s.schema.Get(key); exist {
			if err := k.Validate(v, ctx); err != nil {
				errors = append(errors, err)
			}
		}
	}
}