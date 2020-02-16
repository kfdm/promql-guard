package config

import (
	"os"
	"testing"

	"github.com/prometheus/prometheus/util/testutil"
)

func init() {
	// For finding our test configuration files
	os.Chdir("..")
}

func TestExampleFile(t *testing.T) {
	_, err := LoadFile("guard.yml")
	testutil.Ok(t, err)
}
