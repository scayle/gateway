package resolver

import (
	userService "github.com/scayle/proto/go/user_service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService userService.UserServiceClient
}
