package main

import (
	"github.com/NekoWheel/NekoCAS/db"
	"github.com/NekoWheel/NekoCAS/web"
	log "unknwon.dev/clog/v2"
)

func init() {
	_ = log.NewConsole(100)
}

func main() {
	db.ConnDB()
	db.ConnRedis()
	
	web.Run()
	log.Stop()
}
