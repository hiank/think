package token

import (
	"testing"

	"gotest.tools/v3/assert"
)

var nilToken *Token

func TestGetBuilderOnce(t *testing.T) {
	assert.Equal(t, GetBuilder(), GetBuilder())
}

func TestBuilderGet(t *testing.T) {
	assert.Assert(t, GetBuilder().Get("test") != nilToken)
	assert.Equal(t, GetBuilder().Get("test"), GetBuilder().Get("test"))
}

func TestBuilderFind(t *testing.T) {
	tok, ok := GetBuilder().Find("test1")
	assert.Equal(t, ok, false)
	assert.Equal(t, tok, nilToken)
	tok1 := GetBuilder().Get("test1")
	tok, ok = GetBuilder().Find("test1")
	assert.Equal(t, tok, tok1)
	assert.Equal(t, ok, true)
}
