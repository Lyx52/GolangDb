package server

import (
	"fmt"

	"github.com/Lyx52/GolangDb/schema"
	"github.com/Lyx52/GolangDb/signals"
)

type ServerContext struct {
	CurrentDatabase *schema.Database
	DatabaseStore   *DatabaseStore
	Cancelled       signals.CancelSignal
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

	err, database := context.DatabaseStore.GetDatabase(context.CurrentDatabase.Name)
	if err != nil {
		return err, nil
	}

	err, table := database.CreateTable(name)
	if err != nil {
		return err, nil
	}

	go table.Datastore.HandleDataStoreWrites(context.Cancelled)

	return nil, table
}

func (context *ServerContext) RunDatabases() {
	context.Cancelled = make(signals.CancelSignal)
	for _, database := range context.DatabaseStore.Databases {
		for _, table := range database.Tables {
			go table.Datastore.HandleDataStoreWrites(context.Cancelled)
		}
	}
}
