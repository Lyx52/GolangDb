package schema

import (
	"fmt"
	"math"
)

type ValidateConstraint func(field *TableField, value any) error
type Constraint struct {
	Key       string
	Validator ValidateConstraint
}

func (constr *Constraint) String() string {
	return constr.Key
}

func NewIntegerConstraint(field *TableField) *Constraint {
	return &Constraint{
		Key: fmt.Sprintf("%s_%s_%s::integer", field.table.database.Name, field.table.Name, field.Name),
		Validator: func(field *TableField, value any) error {
			switch converted := value.(type) {
			case float64:
				if converted == math.Trunc(converted) {
					return nil
				}
			}

			return fmt.Errorf("expected integer value")
		},
	}
}

func NewFloatConstraint(field *TableField) *Constraint {
	return &Constraint{
		Key: fmt.Sprintf("%s_%s_%s::float", field.table.database.Name, field.table.Name, field.Name),
		Validator: func(field *TableField, value any) error {
			switch value.(type) {
			case float64:
				return nil
			}

			return fmt.Errorf("expected float value")
		},
	}
}

func NewStringConstraint(field *TableField) *Constraint {
	return &Constraint{
		Key: fmt.Sprintf("%s_%s_%s::string", field.table.database.Name, field.table.Name, field.Name),
		Validator: func(field *TableField, value any) error {
			switch value.(type) {
			case string:
				return nil
			}

			return fmt.Errorf("expected string value")
		},
	}
}
