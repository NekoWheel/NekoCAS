package main

import (
	"github.com/NekoWheel/NekoCAS/internal/conf"
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/web"
	log "unknwon.dev/clog/v2"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	if err = conf.Load(); err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	if err = db.ConnDB(); err != nil {
		log.Fatal("Failed to connect to MySQL database: %v", err)
	}

	if err = db.ConnRedis(); err != nil {
		log.Fatal("Failed to connect to Redis: %v", err)
	}

	web.Run()
}
