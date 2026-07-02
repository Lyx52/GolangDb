package schema

type Database struct {
	Name string
}

func NewDatabase(name string) *Database {
	return &Database{
		Name: name,
	}
}

func (database *Database) String() string {
	return database.Name
}
