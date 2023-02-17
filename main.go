package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/malcolmmaima/maimabank/api"
	db "github.com/malcolmmaima/maimabank/db/sqlc"
)

func main() {
	config, err := util.loadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}