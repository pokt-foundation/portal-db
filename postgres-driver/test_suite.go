package postgresdriver

import (
	"database/sql"
	"errors"
	"fmt"

	// used to embed SQL schema for testing purposes
	_ "embed"

	"github.com/stretchr/testify/suite"
)

var (
	//go:embed sqlc/schema.sql
	schema string

	errConnectingToDB = errors.New("error connecting to test postgres database")
	errInitializingDB = errors.New("error initializing test postgres database")
	errClosingDB      = errors.New("error closing connection to test postgres database")
)

// Connects to a test DB and initializes it with the DB schema
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

type (
	PGDriverTestSuite struct {
		suite.Suite
		connectionString string
	}
)

// SetupSuite runs before each test suite run
func (ts *PGDriverTestSuite) SetupSuite() {
	err := InitializeTestPostgresDB(ts.connectionString)
	ts.NoError(err)
}
