package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/scayle/gateway/graph/generated"
	"github.com/scayle/gateway/graph/model"
	userService "github.com/scayle/proto-go/user_service"
)

func (r *mutationResolver) CreateUser(ctx context.Context, newUser model.NewUser) (*model.CreateUserResponse, error) {
	claims := ctx.Value("claims").(*userService.TokenClaims)

	req := userService.CreateUserRequest{
		Claims:   claims,
		IsAdmin:  newUser.IsAdmin,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
	}

	createReq, err := r.UserService.Create(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	return &model.CreateUserResponse{
		ID: createReq.Id,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, updatedUser model.UpdateUser) (*model.User, error) {
	claims := ctx.Value("claims").(*userService.TokenClaims)

	var isAdmin *wrappers.BoolValue
	if updatedUser.IsAdmin != nil {
		isAdmin = &wrappers.BoolValue{Value: *updatedUser.IsAdmin}
	}
	var username *wrappers.StringValue
	if updatedUser.Username != nil {
		username = &wrappers.StringValue{Value: *updatedUser.Username}
	}
	var email *wrappers.StringValue
	if updatedUser.Email != nil {
		email = &wrappers.StringValue{Value: *updatedUser.Email}
	}
	var password *wrappers.StringValue
	if updatedUser.Password != nil {
		password = &wrappers.StringValue{Value: *updatedUser.Password}
	}

	req := userService.UpdateUserRequest{
		Claims:   claims,
		Id:       updatedUser.ID,
		IsAdmin:  isAdmin,
		Username: username,
		Email:    email,
		Password: password,
	}

	user, err := r.UserService.Update(ctx, &req)

	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	return &model.User{
		ID:    user.Id,
		Name:  user.Username,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
	}, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	claims := ctx.Value("claims").(*userService.TokenClaims)

	req := userService.GetUserRequest{
		Claims: claims,
		Id:     id,
	}

	user, err := r.UserService.Get(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	return &model.User{
		ID:    user.Id,
		Name:  user.Username,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
	}, nil
}

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	claims := ctx.Value("claims").(*userService.TokenClaims)

	req := userService.GetAllUserRequest{
		Claims: claims,
	}

	users, err := r.UserService.GetAll(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	result := make([]*model.User, 0)

	for _, u := range users.Users {
		result = append(result, &model.User{
			ID:    u.Id,
			Name:  u.Username,
			Email: u.Email,
			IsAdmin: u.IsAdmin,
		})
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
