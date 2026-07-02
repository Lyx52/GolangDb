package schema

import "fmt"

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
		Key:       fmt.Sprintf("%s_%s_%s::integer", field.table.database.Name, field.table.Name, field.Name),
		Validator: ValidateIntegerConstraint,
	}
}

func ValidateIntegerConstraint(field *TableField, value any) error {
	return nil
}
