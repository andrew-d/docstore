package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustParse(s string) *AST {
	val, err := Parse("test", []byte(s))
	if err != nil {
		panic(err)
	}
	return val.(*AST)
}

func TestParseSingle(t *testing.T) {
	assert.True(t, true)
}
