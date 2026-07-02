package backing

import (
	"fmt"

	"github.com/Lyx52/GolangDb/schema"
)

type ServerContext struct {
	CurrentDatabase *schema.Database
	DatabaseStore   *DatabaseStore
}

func NewServerContext() *ServerContext {
	return &ServerContext{
		DatabaseStore:   NewDatabaseStore(),
		CurrentDatabase: nil,
	}
}

func (context *ServerContext) CheckDatabaseConnected() error {
	if context.CurrentDatabase == nil {
		return fmt.Errorf("database not connected")
	}

	return nil
}

func (context *ServerContext) SetCurrentDatabase(database *schema.Database) {
	context.CurrentDatabase = database
}

func (context *ServerContext) CreateDatabase(name string) (error, *schema.Database) {
	return context.DatabaseStore.CreateDatabase(name)
}

func (context *ServerContext) CreateTable(name string) (error, *schema.Table) {
	err := context.CheckDatabaseConnected()
	if err != nil {
		return err, nil
	}

	err, tableStore := context.DatabaseStore.GetDatabaseTableStore(context.CurrentDatabase.Name)
	if err != nil {
		return err, nil
	}

	return tableStore.CreateTable(name)
}
