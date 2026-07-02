package objects

import "fmt"

type Table struct {
	Name     string
	Database *Database
}

func NewTable(name string, database *Database) *Table {
	return &Table{
		Name:     name,
		Database: database,
	}
}

func (table *Table) String() string {
	return fmt.Sprintf("%s.%s", table.Database.String(), table.Name)
}
