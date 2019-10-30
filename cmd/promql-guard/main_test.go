package main

import (
	"flag"
	"path/filepath"
	"testing"
)

var testConfig = filepath.Join("..", "..", "example.yaml")

func TestMain(m *testing.M) {
	flag.Parse()
}
