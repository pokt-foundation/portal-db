package driver

import (
	"context"
	"database/sql"

	"github.com/pokt-foundation/portal-db/repository"

	// PQ import is required
	_ "github.com/lib/pq"
)

type (
	PostgresDriver struct {
		Querier
		notification chan *repository.Notification
		listener     Listener
	}
	IPostgresDriver interface {
		Querier
		ReadBlockchains(ctx context.Context) ([]*repository.Blockchain, error)
		WriteBlockchain(ctx context.Context, blockchain *repository.Blockchain) (*repository.Blockchain, error)
		ActivateChain(ctx context.Context, id string, active bool) error
	}
)

/* NewDriver returns PostgresDriver instance from Postgres connection string */
func NewDriver(connectionString string, listener Listener) (*PostgresDriver, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	driver := &PostgresDriver{
		Querier:      New(db),
		notification: make(chan *repository.Notification, 32),
		listener:     listener,
	}

	err = driver.listener.Listen("events")
	if err != nil {
		return nil, err
	}

	go Listen(driver.listener.NotificationChannel(), driver.notification)

	return driver, nil
}

/* NewDriverFromDBInstance returns PostgresDriver instance from sdl.DB instance */
// mostly used for mocking tests
func NewDriverFromDBInstance(db *sql.DB, listener Listener) *PostgresDriver {
	driver := &PostgresDriver{
		Querier:      New(db),
		notification: make(chan *repository.Notification, 32),
		listener:     listener,
	}

	err := driver.listener.Listen("events")
	if err != nil {
		panic(err)
	}

	go Listen(driver.listener.NotificationChannel(), driver.notification)

	return driver
}

/* NotificationChannel returns receiver Notification channel  */
func (d *PostgresDriver) NotificationChannel() <-chan *repository.Notification {
	return d.notification
}

func newSQLNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func newSQLNullInt32(value int32) sql.NullInt32 {
	if value == 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: value,
		Valid: true,
	}
}

func newSQLNullBool(value bool) sql.NullBool {
	return sql.NullBool{
		Bool:  value,
		Valid: true,
	}
}
