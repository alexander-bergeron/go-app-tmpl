package user

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	userpb "github.com/alexander-bergeron/go-app-tmpl/gen/go/proto/user/v1"
	"github.com/alexander-bergeron/go-app-tmpl/internal/repository"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	userpb.UnimplementedUserServiceServer
	q *repository.Queries
}

func NewUserService(db repository.DBTX) *Service {
	return &Service{
		q: repository.New(db),
	}
}

func (s Service) GetUsers(ctx context.Context, _ *emptypb.Empty) (*userpb.GetUsersResponse, error) {
	users, err := s.q.GetAllUsers(ctx)
	if err != nil {
		slog.Error("error querying users", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to read query file: %w", err)
	}

	var pusers []*userpb.User
	for _, u := range users {
		pu := &userpb.User{
			UserId:    u.UserID,
			Username:  u.Username,
			Email:     u.Email,
			Firstname: u.FirstName.String,
			Lastname:  u.LastName.String,
		}
		pusers = append(pusers, pu)
	}

	return &userpb.GetUsersResponse{Users: pusers}, nil
}

func (s Service) CreateUser(ctx context.Context, in *userpb.CreateUserRequest) (*emptypb.Empty, error) {
	newUser := repository.CreateUserParams{
		Username:  in.User.Username,
		Email:     in.User.Email,
		FirstName: sql.NullString{String: in.User.Firstname, Valid: true}, // TODO: check if not null
		LastName:  sql.NullString{String: in.User.Lastname, Valid: true},  // TODO: check if not null
	}

	err := s.q.CreateUser(ctx, newUser)
	if err != nil {
		slog.Error("error creating new user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s Service) StreamUsers(_ *emptypb.Empty, stream userpb.UserService_StreamUsersServer) error {
	// if request.GetIntervalSeconds() == 0 {
	// 	return status.Error(codes.InvalidArgument, "interval must be set")
	// }

	users, err := s.q.GetAllUsers(context.Background())
	if err != nil {
		slog.Error("error querying users", slog.String("error", err.Error()))
		return fmt.Errorf("failed to read query file: %w", err)
	}

	for _, u := range users {
		pu := &userpb.User{
			UserId:    u.UserID,
			Username:  u.Username,
			Email:     u.Email,
			Firstname: u.FirstName.String,
			Lastname:  u.LastName.String,
		}

		err := stream.Send(pu)
		if err != nil {
			return err
		}
	}
	return nil
}
