package postgresdriver

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/pokt-foundation/portal-db/repository"

	// PQ import is required
	_ "github.com/lib/pq"
)

const (
	psqlDateLayout = "2006-01-02T15:04:05.999999"
	idLength       = 24
)

var (
	ErrMissingID = errors.New("missing id")
)

// The PostgresDriver struct satisfies the Source interface which defines all database driver methods
type PostgresDriver struct {
	*Queries
	notification chan *repository.Notification
	listener     Listener
}

/* NewPostgresDriver returns PostgresDriver instance from Postgres connection string */
func NewPostgresDriver(connectionString string, listener Listener) (*PostgresDriver, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	driver := &PostgresDriver{
		Queries:      New(db),
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

/* NewPostgresDriverFromDBInstance returns PostgresDriver instance from sdl.DB instance */
// mostly used for mocking tests
func NewPostgresDriverFromDBInstance(db *sql.DB, listener Listener) *PostgresDriver {
	driver := &PostgresDriver{
		Queries:      New(db),
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

func generateRandomID() (string, error) {
	bytes := make([]byte, idLength/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
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

func newSQLNullTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  value,
		Valid: true,
	}
}

func psqlDateToTime(rawDate string) time.Time {
	date, _ := time.Parse(psqlDateLayout, rawDate)
	return date
}
