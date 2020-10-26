package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go/log"
	"github.com/scayle/gateway/graph/generated"
	"github.com/scayle/gateway/graph/model"
	userService "github.com/scayle/proto/go/user_service"
)

func (r *mutationResolver) CreateUser(ctx context.Context, newUser model.NewUser) (*model.CreateUserResponse, error) {
	claims := ctx.Value("claims").(*userService.TokenClaims)

	u := userService.CreateUserRequest{
		Claims:   claims,
		IsAdmin:  newUser.IsAdmin,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
	}

	createReq, err := r.UserService.Create(ctx, &u)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	return &model.CreateUserResponse{
		ID: createReq.Id,
	}, nil
}

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

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	claims := ctx.Value("claims").(*userService.TokenClaims)

	u := userService.GetUserRequest{
		Claims: claims,
		Id:     id,
	}

	user, err := r.UserService.Get(ctx, &u)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	return &model.User{
		ID:    user.Id,
		Name:  user.Username,
		Email: user.Email,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
