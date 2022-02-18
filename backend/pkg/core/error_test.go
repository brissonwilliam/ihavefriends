package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewErrNotFound(t *testing.T) {
	assert.Equal(t, "Could not find user", NewErrNotFound("user").Error())
}
