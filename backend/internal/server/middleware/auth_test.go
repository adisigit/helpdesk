package middleware

import (
	"context"
	"helpdesk/backend/internal/auth"
	"helpdesk/backend/internal/database/models"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type okOut struct{ Body string }

func setup(t *testing.T) (humatest.TestAPI, *auth.JWT) {
	jwt := auth.NewJWT("0123456789abcdef-secret", time.Hour)
	_, api := humatest.New(t)

	adminOnly := huma.NewGroup(api)
	adminOnly.UseMiddleware(Auth(api, jwt), RequireRole(api, models.RoleAdmin))
	huma.Get(adminOnly, "/admin", func(context.Context, *struct{}) (*okOut, error) {
		return &okOut{Body: "ok"}, nil
	})
	huma.Post(adminOnly, "/admin", func(context.Context, *struct{}) (*okOut, error) {
		return &okOut{Body: "ok"}, nil
	})
	return api, jwt
}

func bearer(t *testing.T, j *auth.JWT, role models.Role) string {
	tok, _, err := j.Generate(uuid.New(), role)
	assert.NoError(t, err)
	return "Authorization: Bearer " + tok
}

func TestAuth_MissingOrInvalidTokenIs401(t *testing.T) {
	api, _ := setup(t)
	assert.Equal(t, http.StatusUnauthorized, api.Get("/admin").Code)
	assert.Equal(t, http.StatusUnauthorized, api.Get("/admin", "Authorization: Bearer sampah").Code)
	assert.Equal(t, http.StatusUnauthorized, api.Get("/admin", "Authorization: Basic abc").Code)
}

func TestRequireRole_RequesterAndAgentGet403OnAdminRoute(t *testing.T) {
	api, jwt := setup(t)
	assert.Equal(t, http.StatusForbidden, api.Get("/admin", bearer(t, jwt, models.RoleRequester)).Code)
	assert.Equal(t, http.StatusForbidden, api.Get("/admin", bearer(t, jwt, models.RoleAgent)).Code)
}

func TestRequireRole_AdminAllowed(t *testing.T) {
	api, jwt := setup(t)
	assert.Equal(t, http.StatusOK, api.Get("/admin", bearer(t, jwt, models.RoleAdmin)).Code)
}

func TestAuth_CookieMutationRequiresCSRFHeader(t *testing.T) {
	api, jwt := setup(t)
	tok, _, _ := jwt.Generate(uuid.New(), models.RoleAdmin)
	cookie := "Cookie: token=" + tok

	assert.Equal(t, http.StatusOK, api.Get("/admin", cookie).Code, "GET via cookie boleh")
	assert.Equal(t, http.StatusForbidden, api.Post("/admin", cookie).Code, "POST via cookie tanpa header ditolak")
	assert.Equal(t, http.StatusOK, api.Post("/admin", cookie, "X-Requested-With: XMLHttpRequest").Code)
	assert.Equal(t, http.StatusOK, api.Post("/admin", bearer(t, jwt, models.RoleAdmin)).Code, "Bearer tidak kena CSRF check")
}
