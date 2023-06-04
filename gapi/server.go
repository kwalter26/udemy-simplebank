package gapi

import (
	"fmt"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/pb"
	"github.com/kwalter26/udemy-simplebank/token"
	"github.com/kwalter26/udemy-simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer Creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maketer: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: maker,
		config:     config,
	}

	return server, nil
}
