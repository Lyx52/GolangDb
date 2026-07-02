package backing

import (
	"fmt"

	"github.com/Lyx52/GolangDb/objects"
)

type DatabaseStore struct {
	Databases            map[string]*objects.Database
	DatabaseToTableStore map[string]*TableStore
}

func NewDatabaseStore() *DatabaseStore {
	return &DatabaseStore{
		Databases:            make(map[string]*objects.Database),
		DatabaseToTableStore: make(map[string]*TableStore),
	}
}

func (store *DatabaseStore) GetDatabase(name string) (error, *objects.Database) {
	database, exists := store.Databases[name]
	if exists {
		return nil, database
	}

	return fmt.Errorf("database %s does not exist", name), nil
}

func (store *DatabaseStore) GetDatabaseTableStore(name string) (error, *TableStore) {
	tableStore, exists := store.DatabaseToTableStore[name]
	if exists {
		return nil, tableStore
	}

	return fmt.Errorf("database %s does not exist", name), nil
}

func (store *DatabaseStore) CreateDatabase(name string) (error, *objects.Database) {
	_, exists := store.Databases[name]
	if exists {
		return fmt.Errorf("database with name %s already exists", name), nil
	}

	store.Databases[name] = objects.NewDatabase(name)
	store.DatabaseToTableStore[name] = NewTableStore(store.Databases[name])
	return nil, store.Databases[name]
}
