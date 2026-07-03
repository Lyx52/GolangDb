package schema

import (
	"fmt"

	"github.com/Lyx52/GolangDb/backing"
)

type Table struct {
	Name      string
	Fields    map[string]*TableField
	Datastore *backing.DataStore
	database  *Database
}

func NewTable(name string, database *Database) *Table {
	return &Table{
		Name:      name,
		Fields:    make(map[string]*TableField),
		database:  database,
		Datastore: backing.NewDataStore(),
	}
}

func (table *Table) String() string {
	return fmt.Sprintf("%s.%s (%v)", table.database.String(), table.Name, len(table.Datastore.Values))
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

	field.Index = len(table.Fields)
	table.Fields[name] = field
	return nil
}
