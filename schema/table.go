package schema

import "fmt"

type Table struct {
	Name     string
	Fields   map[string]*TableField
	database *Database
}

func NewTable(name string, database *Database) *Table {
	return &Table{
		Name:     name,
		Fields:   make(map[string]*TableField),
		database: database,
	}
}

func (table *Table) String() string {
	return fmt.Sprintf("%s.%s", table.database.String(), table.Name)
}

func (table *Table) CreateField(name string, fieldType FieldType, constraints []*Constraint) error {
	_, ok := table.Fields[name]
	if ok {
		return fmt.Errorf("table %s already contains column with name %s", table.Name, name)
	}

	field := NewTableField(name, table)
	implied := fieldType.GetTypeImpliedConstraints(field)
	constraints = append(implied, constraints...)
	var err error
	for _, constraint := range constraints {
		err = field.AddConstraint(constraint)
		if err != nil {
			return err
		}
	}

	table.Fields[name] = field
	return nil
}
