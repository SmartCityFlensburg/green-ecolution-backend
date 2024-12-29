package store_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pgContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	code := 1
	ctx := context.Background()
	defer func() { os.Exit(code) }()
	pgContainer = testutils.SetupPostgres(ctx)
	defer func() { _ = pgContainer.Terminate(ctx) }()

	code = m.Run()
}

func TestStore_NewStore(t *testing.T) {
	t.Run("should create new store", func(t *testing.T) {
		// given
		pool := poolConn(t)

		// when
		s := store.NewStore(pool, sqlc.New(pool))

		// then
		assert.NotNil(t, s)
	})

	t.Run("should panic when querier is nil", func(t *testing.T) {
		// given
		pool := poolConn(t)

		// when
		defer func() {
			assert.NotNil(t, recover())
		}()
		_ = store.NewStore(pool, nil)
	})

	t.Run("should panic when pool is nil", func(t *testing.T) {
		// given
		var pool *pgxpool.Pool

		// when
		defer func() {
			assert.NotNil(t, recover())
		}()
		_ = store.NewStore(pool, sqlc.New(pool))
	})

	t.Run("should panic when pool is nil and querier is nil", func(t *testing.T) {
		// when
		defer func() {
			assert.NotNil(t, recover())
		}()
		_ = store.NewStore(nil, nil)
	})
}

func TestStore_DB(t *testing.T) {
	t.Run("should return db", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s := store.NewStore(pool, sqlc.New(pool))

		// when
		db := s.DB()

		// then
		assert.NotNil(t, db)
	})

	t.Run("should return nil when db is nil", func(t *testing.T) {
		// given
		s := &store.Store{}

		// when
		db := s.DB()

		// then
		assert.Nil(t, db)
	})
}

func TestStore_WithTx(t *testing.T) {
	execMigration(t)
	t.Run("should execute function with transaction", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), func(s *store.Store) error {
			return nil
		})

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when function returns error", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), func(s *store.Store) error {
			return pgx.ErrNoRows
		})

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when function is nil", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), nil)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s := store.NewStore(pool, sqlc.New(pool))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := s.WithTx(ctx, func(s *store.Store) error {
			return nil
		})

		// then
		assert.Error(t, err)
	})

	t.Run("should commit transaction", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s := store.NewStore(pool, sqlc.New(pool))
		sensorID := "sensor-1"
		// when
		err := s.WithTx(context.Background(), func(s *store.Store) error {
			_, _ = s.CreateSensor(context.Background(), &sqlc.CreateSensorParams{
				ID:        sensorID,
				Status:    sqlc.SensorStatusOnline,
				Latitude:  54.801539,
				Longitude: 9.446741,
			})
			return nil
		})

		// then
		assert.NoError(t, err)

		// validate
		got, err := s.GetSensorByID(context.Background(), sensorID)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, sensorID, got.ID)
		assert.Equal(t, sqlc.SensorStatusOnline, got.Status)
		assert.Equal(t, 54.801539, got.Latitude)
		assert.Equal(t, 9.446741, got.Longitude)

		// cleanup
		_ = s.DeleteSensor(context.Background(), sensorID)
	})

	t.Run("should rollback transaction", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s := store.NewStore(pool, sqlc.New(pool))
		sensorID := "sensor-1"

		// when
		err := s.WithTx(context.Background(), func(s *store.Store) error {
			_, _ = s.CreateSensor(context.Background(), &sqlc.CreateSensorParams{
				ID:        sensorID,
				Status:    sqlc.SensorStatusOnline,
				Latitude:  54.801539,
				Longitude: 9.446741,
			})
			return assert.AnError
		})

		// then
		assert.Error(t, err)

		// validate
		got, _ := s.GetSensorByID(context.Background(), sensorID)
		assert.Empty(t, got)
	})
}

func execMigration(t testing.TB) {
	dbURL := dbURL(t)
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("Error while connecting to PostgreSQL: %s", err)
	}

	testutils.ExecMigration(db)
	defer db.Close()
}

func dbURL(t testing.TB) string {
	ctx := context.Background()
	dbURL, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Could not get connection string: %s", err)
	}

	return dbURL
}

func poolConn(t testing.TB) *pgxpool.Pool {
	dbURL := dbURL(t)

	pgxConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		t.Fatalf("Error while parsing PostgreSQL connection string: %s", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		t.Fatalf("Error while connecting to PostgreSQL: %s", err)
	}

	return pool
}

func TestHandleErrorExpanded(t *testing.T) {
	dummyStore := &store.Store{}

	t.Run("HandleError with nil error", func(t *testing.T) {
		// when
		err := dummyStore.HandleError(nil, "Testing nil error")

		// validate
		assert.NoError(t, err)
	})

	t.Run("HandleError with NotFoundError", func(t *testing.T) {
		// when
		err := dummyStore.HandleError(pgx.ErrNoRows, "Fetching entity by ID")

		// validate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "NotFoundError")
		assert.Contains(t, err.Error(), "Fetching entity by ID")
	})

	t.Run("HandleError with DatabaseError", func(t *testing.T) {
		// when
		pgErr := &pgconn.PgError{
			Code:    "23505",
			Message: "Duplicate key value violates unique constraint",
		}
		err := dummyStore.HandleError(pgErr, "Creating new entry")

		// validate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DatabaseError")
		assert.Contains(t, err.Error(), "Creating new entry")
		assert.Contains(t, err.Error(), "23505")
	})

	t.Run("HandleError with UnexpectedError", func(t *testing.T) {
		// when
		unexpectedErr := errors.New("unexpected issue occurred")
		err := dummyStore.HandleError(unexpectedErr, "Unexpected operation")

		// validate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UnexpectedError")
		assert.Contains(t, err.Error(), "Unexpected operation")
	})

	t.Run("HandleError with empty context", func(t *testing.T) {
		// when
		unexpectedErr := errors.New("unexpected issue occurred")
		err := dummyStore.HandleError(unexpectedErr)

		// validate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UnexpectedError")
		assert.Contains(t, err.Error(), "No context provided")
	})

}
