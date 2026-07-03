package schema

import (
	"fmt"
	"strings"
)

type FieldType int

const (
	TypeInteger FieldType = iota
	TypeFloat   FieldType = iota
	TypeVarchar FieldType = iota
)

func (fieldType FieldType) GetTypeImpliedConstraints(field *TableField) []*Constraint {
	switch fieldType {
	case TypeInteger:
		return []*Constraint{NewIntegerConstraint(field)}
	case TypeFloat:
		return []*Constraint{NewFloatConstraint(field)}
	case TypeVarchar:
		return []*Constraint{NewStringConstraint(field)}
	default:
		return make([]*Constraint, 0)
	}
}

func (fieldType FieldType) String() string {
	switch fieldType {
	case TypeInteger:
		return "integer"
	case TypeFloat:
		return "integer"
	case TypeVarchar:
		return "varchar"
	default:
		return "?"
	}
}

func TypeFromString(fieldTypeName string) (error, FieldType) {
	switch fieldTypeName {
	case "integer":
		return nil, TypeInteger
	case "int":
		return nil, TypeInteger
	case "float":
		return nil, TypeFloat
	default:
		if strings.HasPrefix(fieldTypeName, "varchar") {
			return nil, TypeVarchar
		}

		return fmt.Errorf("unknown field type %s", fieldTypeName), -1
	}
}
