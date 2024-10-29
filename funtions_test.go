package templater

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestReduceWhiteSpaces(t *testing.T) {
	assert.Equal(t, "1 2", ReduceWhiteSpaces("  1    2"))
}