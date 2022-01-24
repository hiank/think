package mongo

import (
	"testing"

	mopts "go.mongodb.org/mongo-driver/mongo/options"
	"gotest.tools/v3/assert"
)

func TestKeyConv(t *testing.T) {
	kconv := newKeyConv("@")
	assert.Equal(t, kconv.GetColl(), defaultCollKey)
	assert.Equal(t, kconv.GetDoc(), "")

	kconv = newKeyConv("@11")
	assert.Equal(t, kconv.GetColl(), defaultCollKey)
	assert.Equal(t, kconv.GetDoc(), "11")

	kconv = newKeyConv("11@")
	assert.Equal(t, kconv.GetColl(), defaultCollKey)
	assert.Equal(t, kconv.GetDoc(), "11")

	kconv = newKeyConv("25@gamer")
	assert.Equal(t, kconv.GetColl(), "gamer")
	assert.Equal(t, kconv.GetDoc(), "25")

	kconv = newKeyConv("token")
	assert.Equal(t, kconv.GetColl(), defaultCollKey)
	assert.Equal(t, kconv.GetDoc(), "token")
}

func TestOptions(t *testing.T) {
	dopts := defaultOptions()
	WithDB("test").apply(&dopts)
	assert.Equal(t, dopts.dbName, "test")

	assert.Equal(t, len(dopts.clientOpts), 0)
	WithClientOptions(mopts.Client().ApplyURI("url")).apply(&dopts)
	assert.Equal(t, len(dopts.clientOpts), 1)

	assert.Equal(t, len(dopts.collectionOpts), 0)
	WithCollectionOptions(mopts.Collection().SetRegistry(nil)).apply(&dopts)
	assert.Equal(t, len(dopts.collectionOpts), 1)

	assert.Equal(t, len(dopts.databaseOpts), 0)
	WithDatabaseOptions(mopts.Database().SetReadConcern(nil)).apply(&dopts)
	assert.Equal(t, len(dopts.databaseOpts), 1)

	assert.Equal(t, len(dopts.deleteOpts), 0)
	WithDeleteOptions(mopts.Delete().SetCollation(nil)).apply(&dopts)
	assert.Equal(t, len(dopts.deleteOpts), 1)

	assert.Equal(t, len(dopts.insertOneOpts), 0)
	WithInsertOneOption(mopts.InsertOne().SetBypassDocumentValidation(false)).apply(&dopts)
	assert.Equal(t, len(dopts.insertOneOpts), 1)
}
