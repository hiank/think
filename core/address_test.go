package core_test

import (
	"testing"

	"github.com/hiank/think/core"

	"github.com/stretchr/testify/assert"
)

func TestWithPort(t *testing.T) {

	assert.Equal(t, core.WithPort("192.168.1.22", 1024), "192.168.1.22:1024")
}
