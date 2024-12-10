package user

import (
	"context"
	"fmt"
	"log/slog"

	userpb "github.com/alexander-bergeron/go-app-tmpl/gen/go/proto/user/v1"
	"github.com/alexander-bergeron/go-app-tmpl/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
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
			FirstName: u.FirstName.String,
			LastName:  u.LastName.String,
			Version:   u.Version,
		}
		pusers = append(pusers, pu)
	}

	return &userpb.GetUsersResponse{Users: pusers}, nil
}

func (s Service) CreateUser(ctx context.Context, in *userpb.CreateUserRequest) (*userpb.User, error) {
	newUser := repository.CreateUserParams{
		Username:  in.User.Username,
		Email:     in.User.Email,
		FirstName: pgtype.Text{String: in.User.FirstName, Valid: true},
		LastName:  pgtype.Text{String: in.User.LastName, Valid: true},
		// FirstName: sql.NullString{String: in.User.FirstName, Valid: true}, // TODO: check if not null
		// LastName:  sql.NullString{String: in.User.LastName, Valid: true},  // TODO: check if not null
	}

	user, err := s.q.CreateUser(ctx, newUser)
	if err != nil {
		slog.Error("error creating new user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	return &userpb.User{
		UserId:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Version:   user.Version,
	}, nil
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
			FirstName: u.FirstName.String,
			LastName:  u.LastName.String,
		}

		err := stream.Send(pu)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Service) UpdateUser(ctx context.Context, in *userpb.UpdateUserRequest) (*userpb.User, error) {
	newUser := repository.UpdateUserParams{
		UserID:    in.UserId,
		Username:  in.Username,
		Email:     in.Email,
		FirstName: pgtype.Text{String: in.FirstName, Valid: true},
		LastName:  pgtype.Text{String: in.LastName, Valid: true},
		// FirstName: sql.NullString{String: in.User.FirstName, Valid: true}, // TODO: check if not null
		// LastName:  sql.NullString{String: in.User.LastName, Valid: true},  // TODO: check if not null
	}

	user, err := s.q.UpdateUser(ctx, newUser)
	if err != nil {
		slog.Error("error updating user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &userpb.User{
		UserId:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Version:   user.Version,
	}, nil
}

func (s Service) DeleteUser(ctx context.Context, in *userpb.DeleteUserRequest) (*userpb.User, error) {
	user, err := s.q.DeleteUser(ctx, in.UserId)
	if err != nil {
		slog.Error("error deleting user", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &userpb.User{
		UserId:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Version:   user.Version,
	}, nil
}
