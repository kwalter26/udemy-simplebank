package gapi

import (
	"fmt"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/pb"
	"github.com/kwalter26/udemy-simplebank/token"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/kwalter26/udemy-simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer Creates a new gRPC server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maketer: %w", err)
	}

	server := &Server{
		store:           store,
		tokenMaker:      maker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
