package config

import (
	"github.com/BurntSushi/toml"
)

type DBDriver string

const (
	DriverBadger    DBDriver = "badger"
	DriverCockroach DBDriver = "cockroach"
)

var settings *Settings

func GetSettings() *Settings {
	return settings
}

type Settings struct {
	Debug bool
	User  struct {
		Email          string
		PublicKeyPath  string
		PrivateKeyPath string
	}
	Instance struct {
		URL string
	}
	Daemon struct {
		IP   IPAddress
		Port Port
	}
	Database struct {
		Driver   DBDriver
		Badger   BadgerSettings
		Postgres PostgresSettings
	}
}

type BadgerSettings struct {
	Path      string
	ValuePath string
}

type PostgresSettings struct {
	User        string
	Password    string
	Host        string
	Port        string
	Database    string
	SSLMode     string
	SSLRootCert string
	Options     string
}

func Parse(path string) error {
	if settings == nil {
		settings = &Settings{}
	}
	_, err := toml.DecodeFile(path, settings)
	if err != nil {
		return err
	}
	return nil
}
