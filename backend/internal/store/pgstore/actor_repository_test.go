package pgstore_test

import (
	"testing"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store/pgstore"
	"github.com/stretchr/testify/assert"
)

// TODO: Test with issue
func TestActorCreate(t *testing.T) {
	ta := actor.TestActor()
	db, clear := pgstore.TestStore(dbPath)
	defer clear("actors")

	st := pgstore.New(db)
	err := st.Actor().Create(ta)
	assert.Error(t, err)
	assert.NotNil(t, ta)
}

func TestFindActor(t *testing.T) {
	ta := actor.TestActor()
	db, clear := pgstore.TestStore(dbPath)
	defer clear("actors")

	st := pgstore.New(db)
	err := st.Actor().Create(ta)
	assert.NoError(t, err)
	assert.NotNil(t, ta)

	tmp, err := st.Actor().Find(ta.Id)
	assert.NoError(t, err)
	assert.NotNil(t, tmp)

	assert.Equal(t, tmp, ta)
}

func TestDeleteActor(t *testing.T) {
	ta := actor.TestActor()
	db, clear := pgstore.TestStore(dbPath)
	defer clear("actors")

	st := pgstore.New(db)
	err := st.Actor().Create(ta)
	assert.NoError(t, err)
	assert.NotNil(t, ta)

	id, err := st.Actor().Delete(ta.Id)
	assert.NoError(t, err)
	assert.Equal(t, id, ta.Id)

	tmp, err := st.Actor().Find(ta.Id)
	assert.Equal(t, err, store.ErrRecordNotFound)
	assert.Nil(t, tmp)
}

// TODO: Test func
func TestOverwrightActor(t *testing.T) {
	ta := actor.TestActor()
	db, clear := pgstore.TestStore(dbPath)
	defer clear("actors")

	st := pgstore.New(db)
	err := st.Actor().Create(ta)
	assert.NoError(t, err)
	assert.NotNil(t, ta)

	err = st.Actor().Overwright(ta)
	assert.NoError(t, err)

	err = st.Actor().Overwright(ta)
	assert.Equal(t, err, store.ErrRecordNotFound)
}
