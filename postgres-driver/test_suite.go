package postgresdriver

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	// used to embed SQL schema for testing purposes
	_ "embed"

	"github.com/lib/pq"
	"github.com/pokt-foundation/portal-db/driver"
	"github.com/stretchr/testify/suite"
)

const (
	populateQueryPath = "../testdata/seed_test_db.sql"
	connectionString  = "postgres://postgres:pgpassword@localhost:5432/postgres?sslmode=disable"
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

type (
	PGDriverTestSuite struct {
		suite.Suite
		connectionString string
		driver           driver.Driver
	}
)

// SetupSuite runs before each test suite run
func (ts *PGDriverTestSuite) SetupSuite() {
	err := InitializeTestPostgresDB(ts.connectionString)
	ts.NoError(err)

	err = ts.initPostgresDriver()
	ts.NoError(err)

	err = ts.seedTestDB(populateQueryPath)
	ts.NoError(err)
}

// Initializes a real instance of the Postgres driver that connects to the test Postgres Docker container
func (ts *PGDriverTestSuite) initPostgresDriver() error {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Printf("Problem with listener, error: %s, event type: %d", err.Error(), ev)
		}
	}
	listener := pq.NewListener(ts.connectionString, 10*time.Second, time.Minute, reportProblem)

	driver, err := NewPostgresDriver(ts.connectionString, listener)
	if err != nil {
		return err
	}
	ts.driver = driver

	return nil
}

// Seeds the test Postgres Docker container with test data
func (ts *PGDriverTestSuite) seedTestDB(path string) error {
	queryString, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%w: %s", errReadingSchemaFile, err)
	}
	query := string(queryString)

	db, err := sql.Open("postgres", ts.connectionString)
	if err != nil {
		return fmt.Errorf("%w: %s", errConnectingToDB, err)
	}

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("%w: %s", errSeedingDB, err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("%w: %s", errClosingDB, err)
	}

	return nil
}
