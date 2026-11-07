package router

import (
	"helpdesk/backend/internal/server/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func ProtectedGroup(api huma.API, mws ...middleware.Func) *huma.Group {
	g := huma.NewGroup(api)
	g.UseMiddleware(mws...)
	g.UseSimpleModifier(func(op *huma.Operation) {
		op.Security = []map[string][]string{{"bearer": {}}}
	})
	return g
}
