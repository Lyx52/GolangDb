package backing

import (
	"fmt"

	"github.com/Lyx52/GolangDb/objects"
)

type TableStore struct {
	Tables   map[string]*objects.Table
	Database *objects.Database
}

func NewTableStore(database *objects.Database) *TableStore {
	return &TableStore{
		Tables:   make(map[string]*objects.Table),
		Database: database,
	}
}

func (store *TableStore) GetTable(name string) (error, *objects.Table) {
	table, exists := store.Tables[name]
	if exists {
		return nil, table
	}

	return fmt.Errorf("table %s does not exist", name), nil
}

func (store *TableStore) CreateTable(name string) (error, *objects.Table) {
	_, exists := store.Tables[name]
	if exists {
		return fmt.Errorf("table with name %s already exists", name), nil
	}

	store.Tables[name] = objects.NewTable(name, store.Database)
	return nil, store.Tables[name]
}
