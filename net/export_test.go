package net

import "sync"

var (
	Export_clientm = func(cli *Client) sync.Map {
		return cli.m
	}
)
