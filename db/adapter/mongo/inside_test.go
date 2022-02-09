package mongo

import (
	"testing"

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
