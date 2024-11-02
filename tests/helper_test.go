package tests

import (
	"testing"

	"github.com/AllanM007/simpler-test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestOffset(t *testing.T) {
	page := 2
	limit := 10

	want := 10

	result := helpers.GetOffset(page, limit)

	assert.Equal(t, want, result)
}
