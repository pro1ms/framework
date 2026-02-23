package http

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	scanner := NewScanner()
	assert.NotNil(t, scanner)

	err := os.RemoveAll("./generated")
	assert.NoError(t, err)
	err = scanner.Run("./data/test01", "./generated")
	assert.NoError(t, err)
}
