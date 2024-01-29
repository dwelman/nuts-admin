// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package api

import (
	"github.com/labstack/echo/v4"
)

// Identity An identity object
type Identity struct {
	// Did The DID associated with this identity
	Did string `json:"did"`

	// Name The name of this identity, which is the last path part of a did:web DID.
	// If the DID does not contain paths, or it is not a did:web DID, it will be the same as the DID.
	Name string `json:"name"`
}

// CreateIdentityJSONBody defines parameters for CreateIdentity.
type CreateIdentityJSONBody struct {
	DidQualifier string `json:"did_qualifier"`
}

// CreateIdentityJSONRequestBody defines body for CreateIdentity for application/json ContentType.
type CreateIdentityJSONRequestBody CreateIdentityJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get a list of configured Discovery Services
	// (GET /api/discovery)
	GetDiscoveryServices(ctx echo.Context) error

	// (GET /api/id)
	GetIdentities(ctx echo.Context) error

	// (POST /api/id)
	CreateIdentity(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetDiscoveryServices converts echo context to params.
func (w *ServerInterfaceWrapper) GetDiscoveryServices(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDiscoveryServices(ctx)
	return err
}

// GetIdentities converts echo context to params.
func (w *ServerInterfaceWrapper) GetIdentities(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetIdentities(ctx)
	return err
}

// CreateIdentity converts echo context to params.
func (w *ServerInterfaceWrapper) CreateIdentity(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateIdentity(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/discovery", wrapper.GetDiscoveryServices)
	router.GET(baseURL+"/api/id", wrapper.GetIdentities)
	router.POST(baseURL+"/api/id", wrapper.CreateIdentity)

}
