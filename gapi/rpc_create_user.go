package gapi

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/pb"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/kwalter26/udemy-simplebank/val"
	"github.com/kwalter26/udemy-simplebank/worker"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *Server) CreateUser(context context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		HashedPassword: hashedPassword,
		Email:          req.GetEmail(),
	}

	user, err := s.store.CreateUser(context, arg)
	if err != nil {
		if pgErr, err := err.(*pq.Error); err {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", pgErr)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// Send task to worker
	// TODO: user db transaction
	taskPayload := worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.EmailQueue),
	}

	if err := s.taskDistributor.DistributeTaskSendVerifyEmail(context, &taskPayload, opts...); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task: %s", err)
	}

	response := &pb.CreateUserResponse{
		User: userToPb(user),
	}
	return response, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return violations
}
