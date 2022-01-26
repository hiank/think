package mongo

import (
	mopts "go.mongodb.org/mongo-driver/mongo/options"
)

type options struct {
	dbName         string //mongo Database key
	clientOpts     []*mopts.ClientOptions
	databaseOpts   []*mopts.DatabaseOptions
	collectionOpts []*mopts.CollectionOptions
	insertOneOpts  []*mopts.InsertOneOptions
	findOneOpts    []*mopts.FindOneOptions
	deleteOpts     []*mopts.DeleteOptions
}

type Option interface {
	apply(*options)
}

type funcOption func(*options)

func (fo funcOption) apply(opts *options) {
	fo(opts)
}

func WithDB(name string) Option {
	return funcOption(func(opts *options) {
		opts.dbName = name
	})
}

func WithClientOptions(clientOpts ...*mopts.ClientOptions) Option {
	return funcOption(func(opts *options) {
		opts.clientOpts = clientOpts
	})
}

func WithDatabaseOptions(dbOpts ...*mopts.DatabaseOptions) Option {
	return funcOption(func(opts *options) {
		opts.databaseOpts = dbOpts
	})
}

func WithCollectionOptions(collOpts ...*mopts.CollectionOptions) Option {
	return funcOption(func(opts *options) {
		opts.collectionOpts = collOpts
	})
}

func WithInsertOneOption(insertOneOpts ...*mopts.InsertOneOptions) Option {
	return funcOption(func(opts *options) {
		opts.insertOneOpts = insertOneOpts
	})
}

func WithFindOneOptions(findOneOpts ...*mopts.FindOneOptions) Option {
	return funcOption(func(opts *options) {
		opts.findOneOpts = findOneOpts
	})
}

func WithDeleteOptions(delOpts ...*mopts.DeleteOptions) Option {
	return funcOption(func(opts *options) {
		opts.deleteOpts = delOpts
	})
}
