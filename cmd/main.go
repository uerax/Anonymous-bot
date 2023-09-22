package main

import (
	"flag"
	"log"
	"os"

	"github.com/uerax/Anonymous-bot/config"
	"github.com/uerax/Anonymous-bot/core"
)

var (
	path string
)

func main() {
	setupCmd()
	setupLog()
	setupCfg()
	setupServ()
}

func setupLog() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
}

func setupCmd() {
	flag.StringVar(&path, "c", "../anonymous-bot.yml", "项目的配置文件地址(使用绝对路径) 例: -c /etc/anonymous-bot.yml")
	flag.Parse()
}

func setupCfg() {
	config.Init(path)
}

func setupServ() {
	server := core.NewServer()
	server.Start()
}