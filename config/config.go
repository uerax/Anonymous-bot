package config

import "github.com/uerax/goconf"

func Init(path string) {
	goconf.LoadConfig(path)
}