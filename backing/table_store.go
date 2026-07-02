package backing

import (
	"fmt"

	"github.com/Lyx52/GolangDb/schema"
)

type TableStore struct {
	Tables   map[string]*schema.Table
	Database *schema.Database
}

func NewTableStore(database *schema.Database) *TableStore {
	return &TableStore{
		Tables:   make(map[string]*schema.Table),
		Database: database,
	}
}

func (store *TableStore) GetTable(name string) (error, *schema.Table) {
	table, exists := store.Tables[name]
	if exists {
		return nil, table
	}

	return fmt.Errorf("table %s does not exist", name), nil
}

func (store *TableStore) CreateTable(name string) (error, *schema.Table) {
	_, exists := store.Tables[name]
	if exists {
		return fmt.Errorf("table with name %s already exists", name), nil
	}

	store.Tables[name] = schema.NewTable(name, store.Database)
	return nil, store.Tables[name]
}
