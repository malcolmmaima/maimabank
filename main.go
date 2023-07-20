package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/malcolmmaima/maimabank/api"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
	"github.com/malcolmmaima/maimabank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}