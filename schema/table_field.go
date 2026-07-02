package schema

import (
	"fmt"
	"slices"
)

type TableField struct {
	Name        string
	Type        FieldType
	Constraints []Constraint
	table       *Table
}

func NewTableField(name string, table *Table) *TableField {
	return &TableField{
		Name:  name,
		table: table,
	}
}

func (tt *TableField) AddConstraint(constraint *Constraint) error {
	tt.Constraints = append(tt.Constraints, *constraint)
	return nil
}

func (tt *TableField) RemoveConstraint(constraint *Constraint) error {
	exists := false
	for i, constr := range tt.Constraints {
		if constr.Key == constraint.Key {
			tt.Constraints = slices.Delete(tt.Constraints, i, 1)
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("constraint %s not found in table %s", constraint.String(), tt.table.Name)
	}

	return nil
}
