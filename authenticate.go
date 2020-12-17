package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	userService "github.com/scayle/proto-go/user_service"
)

func validateAndGetClaims(ctx context.Context, client userService.UserServiceClient, token string) (*userService.TokenClaims, error) {
	validationToken := userService.ValidateTokenRequest{
		Token: token,
	}

	authReq, err := client.ValidateToken(ctx, &validationToken)

	if err != nil {
		// client should not get the information why exactly
		return nil, fmt.Errorf("could not authenticate user %w", err)
	}

	return authReq, nil
}

func Authenticator(client userService.UserServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			reqToken := r.Header.Get("Authorization")
			if reqToken == "" {
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			splitToken := strings.Split(reqToken, "Bearer ")
			if len(splitToken) != 2 {
				log.Println("authorization header has a wrong format")
				r = r.WithContext(ctx)
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte{})
				if err != nil {
					log.Printf("could not send %v", err)
				}
				return
			}
			reqToken = splitToken[1]

			claims, err := validateAndGetClaims(ctx, client, reqToken)
			if err != nil {
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			// put it in context
			ctx = context.WithValue(r.Context(), "claims", claims)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
