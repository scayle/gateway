package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go/log"
	"github.com/scayle/gateway/graph/model"
	userService "github.com/scayle/proto-go/user_service"
)

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.LoginUserResponse, error) {
	a := userService.AuthRequest{
		Username: username,
		Password: password,
	}

	authReq, err := r.UserService.Auth(ctx, &a)
	if err != nil {
		log.Error(fmt.Errorf("could not authenticate user %w", err))
		// client should not get the information why exactly
		return nil, fmt.Errorf("could not authenticate user")
	}

	return &model.LoginUserResponse{
		ID:    authReq.Id,
		Token: authReq.Token,
	}, nil
}
