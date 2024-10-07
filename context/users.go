package context

import (
	"blog/models"
	"context"
)

type key string

const (
	userKey key = "user"
)

func WithUser (ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User (ctx context.Context) *models.User {
	user, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}