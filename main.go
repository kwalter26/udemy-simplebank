package main

import (
	"database/sql"
	"github.com/kwalter26/udemy-simplebank/api"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/util"
	_ "github.com/lib/pq"
	_ "github.com/newrelic/go-agent/_integrations/nrpq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".", util.Prod)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
