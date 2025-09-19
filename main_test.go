package main

import (
	"os"
	"testing"

	"backend/test"
)

func TestMain(m *testing.M) {
	test.InitTestDB()    // apply migrations once
	code := m.Run()      // run all tests
	test.CleanupTestDB() // drop tables once
	os.Exit(code)
}
