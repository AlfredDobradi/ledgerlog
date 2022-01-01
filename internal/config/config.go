package config

import (
	"github.com/BurntSushi/toml"
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
		PreferredName  string
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
	Site SiteSettings
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
	SSLMode     string `toml:"ssl_mode"`
	SSLRootCert string `toml:"ssl_root_cert"`
	Cluster     string
}

type SiteSettings struct {
	Title      string
	PublicPath string `toml:"public_path"`
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
