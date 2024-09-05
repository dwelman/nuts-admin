package api

import (
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-admin/discovery"
	"github.com/nuts-foundation/nuts-admin/identity"
	"github.com/nuts-foundation/nuts-admin/issuer"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var _ ServerInterface = (*Wrapper)(nil)

type Wrapper struct {
	Identity      identity.Service
	IssuerService issuer.Service
	Discovery     discovery.Service
}

func (w Wrapper) GetIdentities(ctx echo.Context) error {
	identities, err := w.Identity.List(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, identities)
}

func (w Wrapper) CreateIdentity(ctx echo.Context) error {
	identityRequest := CreateIdentityJSONRequestBody{}
	if err := ctx.Bind(&identityRequest); err != nil {
		return err
	}
	result, err := w.Identity.Create(ctx.Request().Context(), identityRequest.Subject)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, result)
}

func (w Wrapper) GetIdentity(ctx echo.Context, did string) error {
	details, err := w.Identity.Get(ctx.Request().Context(), did)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, details)
}

func (w Wrapper) GetIssuedCredentials(ctx echo.Context, params GetIssuedCredentialsParams) error {
	identities, err := w.Identity.List(ctx.Request().Context())
	if err != nil {
		return err
	}
	result := make([]vc.VerifiableCredential, 0)
	for _, currID := range identities {
		for _, issuerDID := range currID.DIDs {
			credentials, err := w.IssuerService.GetIssuedCredentials(ctx.Request().Context(), issuerDID, strings.Split(params.CredentialTypes, ","))
			if err != nil {
				return err
			}
			result = append(result, credentials...)
		}
	}
	return ctx.JSON(http.StatusOK, result)
}
