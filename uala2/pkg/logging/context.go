package logging

import (
	"context"
	"log/slog"

	"github.com/Bancar/uala-auth-team-go-dependencies/libs/logger"
)

type ctxKey string

const (
	tokenSubKey ctxKey = "token_sub"
)

func WithTokenSub(ctx context.Context, sub string) context.Context {
	return context.WithValue(ctx, tokenSubKey, sub)
}

func ContextAttrs() logger.CtxAttrResolver {
	return logger.CtxAttrResolverFunc(func(ctx context.Context) []slog.Attr {
		attrs := make([]slog.Attr, 0, 1) //nolint:mnd

		if email, ok := ctx.Value(tokenSubKey).(string); ok {
			attrs = append(attrs, slog.String("token_sub", email))
		}

		return attrs
	})
}
