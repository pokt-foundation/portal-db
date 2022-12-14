package postgresdriver

import (
	"database/sql"
	"errors"
	"fmt"

	// used to embed SQL schema for testing purposes
	_ "embed"
)

var (
	//go:embed sqlc/schema.sql
	schema string

	errConnectingToDB    = errors.New("error connecting to test postgres database")
	errInitializingDB    = errors.New("error initializing test postgres database")
	errSeedingDB         = errors.New("error seeding test postgres database")
	errClosingDB         = errors.New("error closing connection to test postgres database")
	errReadingSchemaFile = errors.New("error reading sql schema file")
)

/* Connects to a test DB and initializes it with the DB schema. It embeds the DB schema and is
intended to be used in other repos to initialize the test DB, which is why it's not a method. */
func InitializeTestPostgresDB(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("%w: %s", errConnectingToDB, err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("%w: %s", errInitializingDB, err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("%w: %s", errClosingDB, err)
	}

	return nil
}
