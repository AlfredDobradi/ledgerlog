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
	}
	Instance struct {
		URL string
	}
	Daemon struct {
		IP   IPAddress
		Port Port
	}
	Database struct {
		Path      string
		ValuePath string
	}
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
