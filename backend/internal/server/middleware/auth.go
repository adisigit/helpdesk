package middleware

import (
	"context"
	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/database/models"
	"net/http"
	"slices"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type Func = func(huma.Context, func(huma.Context))

type principalKey struct{}

type Principal struct {
	UserID uuid.UUID
	Role   models.Role
}

func PrincipalFrom(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}

func extractToken(ctx huma.Context) (raw string, viaCookie bool) {
	if t, ok := strings.CutPrefix(ctx.Header("Authorization"), "Bearer "); ok && t != "" {
		return t, false
	}
	r := http.Request{Header: http.Header{"Cookie": {ctx.Header("Cookie")}}}
	if c, err := r.Cookie(auth.CookieName); err == nil && c.Value != "" {
		return c.Value, true
	}
	return "", false
}

func isSafeMethod(m string) bool {
	return m == http.MethodGet || m == http.MethodHead || m == http.MethodOptions
}

func Auth(api huma.API, jwt *auth.JWT) Func {
	return func(ctx huma.Context, next func(huma.Context)) {
		raw, viaCookie := extractToken(ctx)
		if raw == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "authentication required")
			return
		}
		if viaCookie && !isSafeMethod(ctx.Method()) && ctx.Header(auth.CSRFHeader) == "" {
			_ = huma.WriteErr(api, ctx, http.StatusForbidden, "missing "+auth.CSRFHeader+" header")
			return
		}
		claims, err := jwt.Parse(raw)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		uid, err := claims.UserID()
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token subject")
			return
		}
		next(huma.WithValue(ctx, principalKey{}, Principal{UserID: uid, Role: claims.Role}))
	}
}

func RequireRole(api huma.API, allowed ...models.Role) Func {
	return func(ctx huma.Context, next func(huma.Context)) {
		p, ok := PrincipalFrom(ctx.Context())
		if !ok {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "authentication required")
			return
		}
		if !slices.Contains(allowed, p.Role) {
			_ = huma.WriteErr(api, ctx, http.StatusForbidden, "insufficient permissions")
			return
		}
		next(ctx)
	}
}
