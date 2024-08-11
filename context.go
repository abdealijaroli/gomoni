package main

import (
	"context"
)

type AuthContext struct {
	AccountID int
	Email     string
}

type authContextKey struct{}

func NewAuthContext(ctx context.Context, accountID int, email string) context.Context {
	return context.WithValue(ctx, authContextKey{}, &AuthContext{
		AccountID: accountID,
		Email:     email,
	})
}

func GetAuthContext(ctx context.Context) (*AuthContext, bool) {
	auth, ok := ctx.Value(authContextKey{}).(*AuthContext)
	return auth, ok
}
