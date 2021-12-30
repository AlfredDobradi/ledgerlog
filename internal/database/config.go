package database

import "github.com/AlfredDobradi/ledgerlog/internal/config"

var (
	driver config.DBDriver = config.DriverCockroach
)

func Driver() config.DBDriver {
	return driver
}

func SetDriver(newDriver config.DBDriver) {
	driver = newDriver
}
