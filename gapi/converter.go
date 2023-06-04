package gapi

import (
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Convert db.User to pb.User
func userToPb(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
	}
}
