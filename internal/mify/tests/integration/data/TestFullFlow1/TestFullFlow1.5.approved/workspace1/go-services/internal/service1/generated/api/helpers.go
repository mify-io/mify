// THIS FILE IS AUTOGENERATED, DO NOT EDIT
// Generated by mify via OpenAPI Generator

package openapi

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"example.com/namespace/workspace1/go-services/internal/service1/generated/core"
	"github.com/go-chi/chi/v5"
)

// Response return a ServiceResponse struct filled
func Response(code int, body interface{}) ServiceResponse {
	return ServiceResponse{
		Code:    code,
		Headers: nil,
		Body:    body,
	}
}

// ResponseWithHeaders return a ServiceResponse struct filled, including headers
func ResponseWithHeaders(code int, headers map[string][]string, body interface{}) ServiceResponse {
	return ServiceResponse{
		Code:    code,
		Headers: headers,
		Body:    body,
	}
}

// IsZeroValue checks if the val is the zero-ed value.
func IsZeroValue(val interface{}) bool {
	return val == nil || reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())
}

// AssertInterfaceRequired recursively checks each struct in a slice against the callback.
// This method traverse nested slices in a preorder fashion.
func AssertRecurseInterfaceRequired(obj interface{}, callback func(interface{}) error) error {
	return AssertRecurseValueRequired(reflect.ValueOf(obj), callback)
}

// AssertNestedValueRequired checks each struct in the nested slice against the callback.
// This method traverse nested slices in a preorder fashion.
func AssertRecurseValueRequired(value reflect.Value, callback func(interface{}) error) error {
	switch value.Kind() {
	// If it is a struct we check using callback
	case reflect.Struct:
		if err := callback(value.Interface()); err != nil {
			return err
		}

	// If it is a slice we continue recursion
	case reflect.Slice:
		for i := 0; i < value.Len(); i += 1 {
			if err := AssertRecurseValueRequired(value.Index(i), callback); err != nil {
				return err
			}
		}
	}
	return nil
}

// wrapper to prevent generating controllers with unused imports
func getBodyDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

// wrapper to prevent generating controllers with unused imports
func getURLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

// extract mify request context builder from http context
func mustGetContextBuilder(r *http.Request) *core.MifyRequestContextBuilder {
	ctxBuilder := r.Context().Value(MifyContextField).(*core.MifyRequestContextBuilder)
	if ctxBuilder == nil {
		panic("Context builder wasn't found")
	}
	return ctxBuilder
}
