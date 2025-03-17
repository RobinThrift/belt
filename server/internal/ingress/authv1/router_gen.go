//lint:file-ignore ST1005 Ignore because generated code
//lint:file-ignore SA1029 Ignore because generated code
//go:build go1.22

// Package authv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package authv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.robinthrift.com/belt/internal/x/httperrors"
)

const (
	TokenBearerAuthScopes = "tokenBearerAuth.Scopes"
)

// Defines values for AuthTokenRequestPasswordGrantGrantType.
const (
	Password AuthTokenRequestPasswordGrantGrantType = "password"
)

// Defines values for AuthTokenRequestRefreshTokenGrantGrantType.
const (
	RefreshToken AuthTokenRequestRefreshTokenGrantGrantType = "refresh_token"
)

// AccountKey An account's public key.
type AccountKey struct {
	Data []byte `json:"data"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// AuthToken AuthToken is a pair of access and refresh tokens.
type AuthToken struct {
	AccessToken      string    `json:"accessToken"`
	ExpiresAt        time.Time `json:"expiresAt"`
	RefreshExpiresAt time.Time `json:"refreshExpiresAt"`
	RefreshToken     string    `json:"refreshToken"`
}

// AuthTokenRequest Auth token request union
type AuthTokenRequest struct {
	union json.RawMessage
}

// AuthTokenRequestPasswordGrant Request a new auth token using username and password.
type AuthTokenRequestPasswordGrant struct {
	GrantType AuthTokenRequestPasswordGrantGrantType `json:"grant_type"`
	Password  string                                 `json:"password"`
	Username  string                                 `json:"username"`
}

// AuthTokenRequestPasswordGrantGrantType defines model for AuthTokenRequestPasswordGrant.GrantType.
type AuthTokenRequestPasswordGrantGrantType string

// AuthTokenRequestRefreshTokenGrant Request a new auth token using a refresh token.
type AuthTokenRequestRefreshTokenGrant struct {
	GrantType    AuthTokenRequestRefreshTokenGrantGrantType `json:"grant_type"`
	RefreshToken string                                     `json:"refresh_token"`
}

// AuthTokenRequestRefreshTokenGrantGrantType defines model for AuthTokenRequestRefreshTokenGrant.GrantType.
type AuthTokenRequestRefreshTokenGrantGrantType string

// Error Follows RFC7807 (https://datatracker.ietf.org/doc/html/rfc7807)
type Error = httperrors.Error

// ErrorBadRequest Follows RFC7807 (https://datatracker.ietf.org/doc/html/rfc7807)
type ErrorBadRequest = Error

// ErrorNotFound Follows RFC7807 (https://datatracker.ietf.org/doc/html/rfc7807)
type ErrorNotFound = Error

// ErrorOther Follows RFC7807 (https://datatracker.ietf.org/doc/html/rfc7807)
type ErrorOther = Error

// ErrorUnauthorized Follows RFC7807 (https://datatracker.ietf.org/doc/html/rfc7807)
type ErrorUnauthorized = Error

// AddAccountKeyRequest The public key name, type and data.
type AddAccountKeyRequest struct {
	Data []byte `json:"data"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ChangePasswordRequest Data for the current and new password
type ChangePasswordRequest struct {
	CurrentPassword   string `json:"currentPassword"`
	NewPassword       string `json:"newPassword"`
	NewPasswordRepeat string `json:"newPasswordRepeat"`
	Username          string `json:"username"`
}

// ChangePasswordJSONBody defines parameters for ChangePassword.
type ChangePasswordJSONBody struct {
	CurrentPassword   string `json:"currentPassword"`
	NewPassword       string `json:"newPassword"`
	NewPasswordRepeat string `json:"newPasswordRepeat"`
	Username          string `json:"username"`
}

// CheckAccessParams defines parameters for CheckAccess.
type CheckAccessParams struct {
	// Authorization Access bearer token.
	Authorization string `json:"Authorization"`
}

// AddAccountKeyJSONBody defines parameters for AddAccountKey.
type AddAccountKeyJSONBody struct {
	Data []byte `json:"data"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ChangePasswordJSONRequestBody defines body for ChangePassword for application/json ContentType.
type ChangePasswordJSONRequestBody ChangePasswordJSONBody

// AddAccountKeyJSONRequestBody defines body for AddAccountKey for application/json ContentType.
type AddAccountKeyJSONRequestBody AddAccountKeyJSONBody

// RequestAuthTokenJSONRequestBody defines body for RequestAuthToken for application/json ContentType.
type RequestAuthTokenJSONRequestBody = AuthTokenRequest

// AsAuthTokenRequestPasswordGrant returns the union data inside the AuthTokenRequest as a AuthTokenRequestPasswordGrant
func (t AuthTokenRequest) AsAuthTokenRequestPasswordGrant() (AuthTokenRequestPasswordGrant, error) {
	var body AuthTokenRequestPasswordGrant
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAuthTokenRequestPasswordGrant overwrites any union data inside the AuthTokenRequest as the provided AuthTokenRequestPasswordGrant
func (t *AuthTokenRequest) FromAuthTokenRequestPasswordGrant(v AuthTokenRequestPasswordGrant) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAuthTokenRequestPasswordGrant performs a merge with any union data inside the AuthTokenRequest, using the provided AuthTokenRequestPasswordGrant
func (t *AuthTokenRequest) MergeAuthTokenRequestPasswordGrant(v AuthTokenRequestPasswordGrant) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsAuthTokenRequestRefreshTokenGrant returns the union data inside the AuthTokenRequest as a AuthTokenRequestRefreshTokenGrant
func (t AuthTokenRequest) AsAuthTokenRequestRefreshTokenGrant() (AuthTokenRequestRefreshTokenGrant, error) {
	var body AuthTokenRequestRefreshTokenGrant
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAuthTokenRequestRefreshTokenGrant overwrites any union data inside the AuthTokenRequest as the provided AuthTokenRequestRefreshTokenGrant
func (t *AuthTokenRequest) FromAuthTokenRequestRefreshTokenGrant(v AuthTokenRequestRefreshTokenGrant) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAuthTokenRequestRefreshTokenGrant performs a merge with any union data inside the AuthTokenRequest, using the provided AuthTokenRequestRefreshTokenGrant
func (t *AuthTokenRequest) MergeAuthTokenRequestRefreshTokenGrant(v AuthTokenRequestRefreshTokenGrant) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t AuthTokenRequest) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *AuthTokenRequest) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Change acocunt password.
	// (POST /change-password)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	// Check if the provided access token is valid.
	// (GET /check-access)
	CheckAccess(w http.ResponseWriter, r *http.Request, params CheckAccessParams)
	// Add a new account key.
	// (POST /keys)
	AddAccountKey(w http.ResponseWriter, r *http.Request)
	// Get a public key by name.
	// (GET /keys/{name})
	GetAccountKey(w http.ResponseWriter, r *http.Request, name string)
	// Request a new AuthToken pair.
	// (POST /token)
	RequestAuthToken(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ChangePassword operation middleware
func (siw *ServerInterfaceWrapper) ChangePassword(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, TokenBearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ChangePassword(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CheckAccess operation middleware
func (siw *ServerInterfaceWrapper) CheckAccess(w http.ResponseWriter, r *http.Request) {

	var err error

	ctx := r.Context()

	ctx = context.WithValue(ctx, TokenBearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	// Parameter object where we will unmarshal all parameters from the context
	var params CheckAccessParams

	headers := r.Header

	// ------------- Required header parameter "Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Authorization", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Authorization", valueList[0], &Authorization, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Authorization", Err: err})
			return
		}

		params.Authorization = Authorization

	} else {
		err := fmt.Errorf("Header parameter Authorization is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Authorization", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CheckAccess(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// AddAccountKey operation middleware
func (siw *ServerInterfaceWrapper) AddAccountKey(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, TokenBearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddAccountKey(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetAccountKey operation middleware
func (siw *ServerInterfaceWrapper) GetAccountKey(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithOptions("simple", "name", r.PathValue("name"), &name, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, TokenBearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAccountKey(w, r, name)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// RequestAuthToken operation middleware
func (siw *ServerInterfaceWrapper) RequestAuthToken(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, TokenBearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RequestAuthToken(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/change-password", wrapper.ChangePassword)
	m.HandleFunc("GET "+options.BaseURL+"/check-access", wrapper.CheckAccess)
	m.HandleFunc("POST "+options.BaseURL+"/keys", wrapper.AddAccountKey)
	m.HandleFunc("GET "+options.BaseURL+"/keys/{name}", wrapper.GetAccountKey)
	m.HandleFunc("POST "+options.BaseURL+"/token", wrapper.RequestAuthToken)

	return m
}

type ErrorBadRequestJSONResponse Error

type ErrorNotFoundJSONResponse Error

type ErrorOtherJSONResponse Error

type ErrorUnauthorizedJSONResponse Error

type ChangePasswordRequestObject struct {
	Body *ChangePasswordJSONRequestBody
}

type ChangePasswordResponseObject interface {
	VisitChangePasswordResponse(w http.ResponseWriter) error
}

type ChangePassword204Response struct {
}

func (response ChangePassword204Response) VisitChangePasswordResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type ChangePassword400JSONResponse struct{ ErrorBadRequestJSONResponse }

func (response ChangePassword400JSONResponse) VisitChangePasswordResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type ChangePassword401JSONResponse struct{ ErrorUnauthorizedJSONResponse }

func (response ChangePassword401JSONResponse) VisitChangePasswordResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type ChangePassword404JSONResponse struct{ ErrorNotFoundJSONResponse }

func (response ChangePassword404JSONResponse) VisitChangePasswordResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type ChangePassworddefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response ChangePassworddefaultJSONResponse) VisitChangePasswordResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type CheckAccessRequestObject struct {
	Params CheckAccessParams
}

type CheckAccessResponseObject interface {
	VisitCheckAccessResponse(w http.ResponseWriter) error
}

type CheckAccess204Response struct {
}

func (response CheckAccess204Response) VisitCheckAccessResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type CheckAccess401JSONResponse struct{ ErrorUnauthorizedJSONResponse }

func (response CheckAccess401JSONResponse) VisitCheckAccessResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type CheckAccessdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response CheckAccessdefaultJSONResponse) VisitCheckAccessResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type AddAccountKeyRequestObject struct {
	Body *AddAccountKeyJSONRequestBody
}

type AddAccountKeyResponseObject interface {
	VisitAddAccountKeyResponse(w http.ResponseWriter) error
}

type AddAccountKey201Response struct {
}

func (response AddAccountKey201Response) VisitAddAccountKeyResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type AddAccountKey400JSONResponse struct{ ErrorBadRequestJSONResponse }

func (response AddAccountKey400JSONResponse) VisitAddAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type AddAccountKey401JSONResponse struct{ ErrorUnauthorizedJSONResponse }

func (response AddAccountKey401JSONResponse) VisitAddAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type AddAccountKey404JSONResponse struct{ ErrorNotFoundJSONResponse }

func (response AddAccountKey404JSONResponse) VisitAddAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type AddAccountKeydefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response AddAccountKeydefaultJSONResponse) VisitAddAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetAccountKeyRequestObject struct {
	Name string `json:"name"`
}

type GetAccountKeyResponseObject interface {
	VisitGetAccountKeyResponse(w http.ResponseWriter) error
}

type GetAccountKey200JSONResponse AccountKey

func (response GetAccountKey200JSONResponse) VisitGetAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetAccountKey400JSONResponse struct{ ErrorBadRequestJSONResponse }

func (response GetAccountKey400JSONResponse) VisitGetAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetAccountKey401JSONResponse struct{ ErrorUnauthorizedJSONResponse }

func (response GetAccountKey401JSONResponse) VisitGetAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetAccountKey404JSONResponse struct{ ErrorNotFoundJSONResponse }

func (response GetAccountKey404JSONResponse) VisitGetAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetAccountKeydefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response GetAccountKeydefaultJSONResponse) VisitGetAccountKeyResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type RequestAuthTokenRequestObject struct {
	Body *RequestAuthTokenJSONRequestBody
}

type RequestAuthTokenResponseObject interface {
	VisitRequestAuthTokenResponse(w http.ResponseWriter) error
}

type RequestAuthToken201JSONResponse AuthToken

func (response RequestAuthToken201JSONResponse) VisitRequestAuthTokenResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type RequestAuthToken204ResponseHeaders struct {
	Location string
}

type RequestAuthToken204Response struct {
	Headers RequestAuthToken204ResponseHeaders
}

func (response RequestAuthToken204Response) VisitRequestAuthTokenResponse(w http.ResponseWriter) error {
	w.Header().Set("Location", fmt.Sprint(response.Headers.Location))
	w.WriteHeader(204)
	return nil
}

type RequestAuthToken400JSONResponse struct{ ErrorBadRequestJSONResponse }

func (response RequestAuthToken400JSONResponse) VisitRequestAuthTokenResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type RequestAuthToken401JSONResponse struct{ ErrorUnauthorizedJSONResponse }

func (response RequestAuthToken401JSONResponse) VisitRequestAuthTokenResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type RequestAuthToken404JSONResponse struct{ ErrorNotFoundJSONResponse }

func (response RequestAuthToken404JSONResponse) VisitRequestAuthTokenResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type RequestAuthTokendefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response RequestAuthTokendefaultJSONResponse) VisitRequestAuthTokenResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Change acocunt password.
	// (POST /change-password)
	ChangePassword(ctx context.Context, request ChangePasswordRequestObject) (ChangePasswordResponseObject, error)
	// Check if the provided access token is valid.
	// (GET /check-access)
	CheckAccess(ctx context.Context, request CheckAccessRequestObject) (CheckAccessResponseObject, error)
	// Add a new account key.
	// (POST /keys)
	AddAccountKey(ctx context.Context, request AddAccountKeyRequestObject) (AddAccountKeyResponseObject, error)
	// Get a public key by name.
	// (GET /keys/{name})
	GetAccountKey(ctx context.Context, request GetAccountKeyRequestObject) (GetAccountKeyResponseObject, error)
	// Request a new AuthToken pair.
	// (POST /token)
	RequestAuthToken(ctx context.Context, request RequestAuthTokenRequestObject) (RequestAuthTokenResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// ChangePassword operation middleware
func (sh *strictHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var request ChangePasswordRequestObject

	var body ChangePasswordJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ChangePassword(ctx, request.(ChangePasswordRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ChangePassword")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ChangePasswordResponseObject); ok {
		if err := validResponse.VisitChangePasswordResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CheckAccess operation middleware
func (sh *strictHandler) CheckAccess(w http.ResponseWriter, r *http.Request, params CheckAccessParams) {
	var request CheckAccessRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CheckAccess(ctx, request.(CheckAccessRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CheckAccess")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CheckAccessResponseObject); ok {
		if err := validResponse.VisitCheckAccessResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// AddAccountKey operation middleware
func (sh *strictHandler) AddAccountKey(w http.ResponseWriter, r *http.Request) {
	var request AddAccountKeyRequestObject

	var body AddAccountKeyJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.AddAccountKey(ctx, request.(AddAccountKeyRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AddAccountKey")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(AddAccountKeyResponseObject); ok {
		if err := validResponse.VisitAddAccountKeyResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetAccountKey operation middleware
func (sh *strictHandler) GetAccountKey(w http.ResponseWriter, r *http.Request, name string) {
	var request GetAccountKeyRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetAccountKey(ctx, request.(GetAccountKeyRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetAccountKey")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetAccountKeyResponseObject); ok {
		if err := validResponse.VisitGetAccountKeyResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// RequestAuthToken operation middleware
func (sh *strictHandler) RequestAuthToken(w http.ResponseWriter, r *http.Request) {
	var request RequestAuthTokenRequestObject

	var body RequestAuthTokenJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.RequestAuthToken(ctx, request.(RequestAuthTokenRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "RequestAuthToken")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(RequestAuthTokenResponseObject); ok {
		if err := validResponse.VisitRequestAuthTokenResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
