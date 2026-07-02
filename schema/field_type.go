package schema

import "fmt"

type FieldType int

const (
	TypeInteger FieldType = iota
	TypeVarchar FieldType = iota
)

func (fieldType FieldType) GetTypeImpliedConstraints(field *TableField) []*Constraint {
	switch fieldType {
	case TypeInteger:
		return []*Constraint{NewIntegerConstraint(field)}
	default:
		return make([]*Constraint, 0)
	}
}

func (fieldType FieldType) String() string {
	switch fieldType {
	case TypeInteger:
		return "integer"
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
	default:
		return fmt.Errorf("unknown field type %s", fieldTypeName), -1
	}
}
