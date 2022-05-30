package net

import "sync"

var (
	Export_clientsetm = func(cs Clientset) sync.Map {
		return cs.(*clientset).m
	}
)
