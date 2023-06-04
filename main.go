package main

import (
	"database/sql"
	"github.com/kwalter26/udemy-simplebank/api"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/gapi"
	"github.com/kwalter26/udemy-simplebank/pb"
	"github.com/kwalter26/udemy-simplebank/util"
	_ "github.com/lib/pq"
	_ "github.com/newrelic/go-agent/_integrations/nrpq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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
	runGRPCServer(config, store)
	runGINServer(config, store)
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create grpc server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Printf("starting gRPC server on %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func runGINServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create gin server:", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
