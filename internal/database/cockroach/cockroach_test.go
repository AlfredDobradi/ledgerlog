package cockroach

import (
	"testing"

	"github.com/AlfredDobradi/ledgerlog/internal/config"
)

func TestBuildConnectionString(t *testing.T) {
	cfg := config.Settings{
		Database: struct {
			Driver   config.DBDriver
			Badger   config.BadgerSettings
			Postgres config.PostgresSettings
		}{
			Driver: config.DriverCockroach,
			Postgres: config.PostgresSettings{
				User:        "test",
				Password:    "pass",
				Host:        "1.1.1.1",
				Port:        "33333",
				Database:    "testingdb",
				SSLMode:     "verify-full",
				SSLRootCert: "asd.crt",
				Options:     "--test%3Dtrue",
			},
		},
	}
	expectedConnectionString := "postgresql://test:pass@1.1.1.1:33333/testingdb?options=--test%253Dtrue&sslmode=verify-full&sslrootcert=asd.crt"

	if actual, expected := buildConnectionString(cfg.Database.Postgres), expectedConnectionString; actual != expected {
		t.Fatalf("Fail: Connection strings don't match.\nExpected: %s\nGot: %s\n", expected, actual)
	}

	t.Log("Pass")
}
