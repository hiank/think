package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestWithPort(t *testing.T) {

	assert.Equal(t, WithPort("192.168.1.22", 1024), "192.168.1.22:1024")
}