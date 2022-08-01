package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	// For finding our test configuration files
	os.Chdir("..")
}

func TestExampleFile(t *testing.T) {
	_, err := LoadFile("guard.yml")
	require.NoError(t, err)
}
