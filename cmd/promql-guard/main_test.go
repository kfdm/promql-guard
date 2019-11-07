package main

import (
	"flag"
	"path/filepath"
	"testing"
)

var testConfig = filepath.Join("..", "..", "guard.yml")

func TestMain(m *testing.M) {
	flag.Parse()
}
