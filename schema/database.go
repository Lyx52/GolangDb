package schema

import (
	"fmt"
)

type Database struct {
	Name   string
	Tables map[string]*Table
}

func NewDatabase(name string) *Database {
	return &Database{
		Name:   name,
		Tables: make(map[string]*Table),
	}
}

func (database *Database) String() string {
	return database.Name
}

func (database *Database) GetTable(name string) (error, *Table) {
	table, exists := database.Tables[name]
	if exists {
		return nil, table
	}

	return fmt.Errorf("table %s does not exist", name), nil
}

func (database *Database) CreateTable(name string) (error, *Table) {
	_, exists := database.Tables[name]
	if exists {
		return fmt.Errorf("table with name %s already exists", name), nil
	}

	database.Tables[name] = NewTable(name, database)
	return nil, database.Tables[name]
}
