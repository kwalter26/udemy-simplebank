package gapi

import (
	"context"
	"fmt"
	"github.com/kwalter26/udemy-simplebank/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (s *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing context metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := fields[0]
	if strings.ToLower(authType) != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type %s", authType)
	}

	authToken := fields[1]
	payload, err := s.tokenMaker.VerifyToken(authToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return payload, nil
}
