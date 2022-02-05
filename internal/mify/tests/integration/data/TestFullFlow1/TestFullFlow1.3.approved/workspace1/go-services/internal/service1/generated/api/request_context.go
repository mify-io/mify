// THIS FILE IS AUTOGENERATED, DO NOT EDIT
// Generated by mify via OpenAPI Generator

package openapi

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"example.com/namespace/workspace1/go-services/internal/service1/generated/core"
)

type ctxKeyMifyContext int

const MifyContextField ctxKeyMifyContext = 0

func RequestContext(sc *core.MifyServiceContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			mifyCtxBuilder := core.NewMifyRequestContextBuilder(sc)
			mifyCtxBuilder.SetRequestID(middleware.GetReqID(r.Context()))
			mifyCtxBuilder.SetProtocol(r.Proto)
			mifyCtxBuilder.SetURLPath(r.URL.Path)
			ctx := context.WithValue(r.Context(), MifyContextField, mifyCtxBuilder)

			next.ServeHTTP(ww, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
