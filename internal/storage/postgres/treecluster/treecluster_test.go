package treecluster

import (
	"context"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	mappers TreeClusterMappers
	suite   *testutils.PostgresTestSuite
)

func TestMain(m *testing.M) {
	code := 1
	ctx := context.Background()
	defer func() { os.Exit(code) }()
	suite = testutils.SetupPostgresTestSuite(ctx)
	mappers = NewTreeClusterRepositoryMappers(
		&generated.InternalTreeClusterRepoMapperImpl{},
		&generated.InternalSensorRepoMapperImpl{},
		&generated.InternalRegionRepoMapperImpl{},
		&generated.InternalTreeRepoMapperImpl{},
	)
	defer suite.Terminate(ctx)

	code = m.Run()
}

func TestTreeClusterRepository_Delete(t *testing.T) {
	t.Run("Delete tree cluster should delete tree cluster", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.ExecQuery(t, "INSERT INTO tree_clusters (id, name, moisture_level, address, description) VALUES (1, 'test', 0.5, '', '')")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 1)
		got, errGot := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.Error(t, errGot)
		assert.Nil(t, got)
	})

	t.Run("Delete tree cluster with linked tree should return error", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 1)

		// then
		assert.Error(t, err)
	})

	t.Run("Delete tree cluster with non-existing id should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("Delete tree cluster with negative id should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("Delete tree cluster with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.Delete(ctx, 1)

		// then
		assert.Error(t, err)
	})
}

func TestTreeClusterRepository_Archived(t *testing.T) {
  t.Run("Archive tree cluster should archive tree cluster", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)

    // when
    err := r.Archive(context.Background(), 1)
    got, errGot := r.GetByID(context.Background(), 1)

    // then
    assert.NoError(t, err)
    assert.NoError(t, errGot)
    assert.NotNil(t, got)
    assert.True(t, got.Archived)
  })

  t.Run("Archive tree cluster with non-existing id should return error", func(t *testing.T) {
    // given
    r := NewTreeClusterRepository(suite.Store, mappers)

    // when
    err := r.Archive(context.Background(), 99)

    // then
    assert.Error(t, err)
  })

  t.Run("Archive tree cluster with negative id should return error", func(t *testing.T) {
    // given
    r := NewTreeClusterRepository(suite.Store, mappers)

    // when
    err := r.Archive(context.Background(), -1)

    // then
    assert.Error(t, err)
  })

  t.Run("Archive tree cluster that is already archived should return not return error", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)
    err := r.Archive(context.Background(), 1)
    assert.NoError(t, err)

    // when
    gotBefore, errGotBefore := r.GetByID(context.Background(), 1)
    err = r.Archive(context.Background(), 1)
    gotAfter, errGotAfter := r.GetByID(context.Background(), 1)

    // then
    assert.NoError(t, err)
    assert.NoError(t, errGotBefore)
    assert.NoError(t, errGotAfter)
    assert.NotNil(t, gotBefore)
    assert.NotNil(t, gotAfter)
    assert.True(t, gotBefore.Archived)
    assert.True(t, gotAfter.Archived)
  })

  t.Run("Archive tree cluster with context canceled exceeded should return error", func(t *testing.T) {
    // given
    r := NewTreeClusterRepository(suite.Store, mappers)
    ctx, cancel := context.WithCancel(context.Background())
    cancel()

    // when
    err := r.Archive(ctx, 1)

    // then
    assert.Error(t, err)
  })
}
