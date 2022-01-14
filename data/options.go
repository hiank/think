package dset

import (
	"github.com/hiank/think/data/db"
)

type options struct {
	buildGamer BuildGamer
	memoryDB   db.IClient
	diskDB     db.IClient
}

type Option interface {
	apply(*options)
}

type funcOption func(*options)

func (fo funcOption) apply(opts *options) {
	fo(opts)
}

//WithMemoryDB memory database client
//NOTE: the memory must be set
func WithMemoryDB(cli db.IClient) Option {
	return funcOption(func(opts *options) {
		opts.memoryDB = cli
	})
}

//WithDiskDB disk database client
func WithDiskDB(cli db.IClient) Option {
	return funcOption(func(opts *options) {
		opts.diskDB = cli
	})
}

//WithGamerBuilder gamer object builder
//the builder given the method to build a gamer object
//liteHub use HGet Scan to unmarshal the value to the gamer object
//and return the value
func WithGamerBuilder(build BuildGamer) Option {
	return funcOption(func(opts *options) {
		opts.buildGamer = build
	})
}

// func With
