package googleit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseTor(t *testing.T) {
	assert.Nil(t, useTor())
}
