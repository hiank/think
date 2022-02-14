package db

import "github.com/hiank/think/run"

const (
	//not found doc in database
	ErrNotFound = run.Err("db: not found in database")
	//invalid param for BytesCoder. must be PB/GOB/JSON
	ErrInvalidBCParam = run.Err("db: invalid BytesCoder param: only support PB GOB JSON now")
	//
	ErrInvalidKey    = run.Err("db: invalid key (should form [`kt`@KT]`baseKey`)")
	ErrUnimplemented = run.Err("db: unimplemented interface")
)
