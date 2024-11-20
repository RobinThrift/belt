package ui

import (
	"fmt"
	"net/http"

	"github.com/RobinThrift/belt/internal/auth"
	"github.com/RobinThrift/belt/internal/control"
	"github.com/RobinThrift/belt/internal/server/session"
	"github.com/RobinThrift/belt/ui"
	"github.com/gorilla/csrf"
)

type router struct {
	authCtrl *control.AuthController
	config   Config
}

type Config struct {
	CSRFSecret       []byte
	UseSecureCookies bool
	BasePath         string
}

func NewRouter(config Config, mux *http.ServeMux, authCtrl *control.AuthController) {
	router := &router{authCtrl: authCtrl, config: config}

	csrfProtectionMiddleware := csrf.Protect(
		config.CSRFSecret,
		csrf.Secure(config.UseSecureCookies),
		csrf.FieldName("belt.csrf.token"),
		csrf.CookieName("belt_csrf_token"),
		csrf.Path(config.BasePath),
		csrf.ErrorHandler(http.HandlerFunc(router.csrfErrorHandler)),
	)

	mux.Handle("/", csrfProtectionMiddleware(router.ensureLoggedIn(router.handlerFuncWithErr(func(w http.ResponseWriter, r *http.Request) error {
		err := router.renderUI(w, r, ui.PageData{Title: "Belt"})
		if err != nil {
			router.renderErrorPage(w, r, err)
		}
		return nil
	}))))

	mux.Handle("/assets/", ui.Assets(config.BasePath+"assets/"))

	mux.Handle("GET /login", csrfProtectionMiddleware(router.handlerFuncWithErr((router.getLogin))))
	mux.Handle("POST /login", csrfProtectionMiddleware(router.handlerFuncWithErr(router.postLogin)))

	mux.Handle("GET /auth/change_password", csrfProtectionMiddleware(router.handlerFuncWithErr((router.getChangePassword))))
	mux.Handle("POST /auth/change_password", csrfProtectionMiddleware(router.handlerFuncWithErr(router.postChangePassord)))

	mux.Handle("GET /logout", router.handlerFuncWithErr(router.getLogout))
}

func (router *router) handlerFuncWithErr(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			router.renderErrorPage(w, r, err)
		}
	}
}

func (router *router) renderUI(w http.ResponseWriter, r *http.Request, data ui.PageData) error {
	data.BaseURL = router.config.BasePath
	data.AssetURL = joinPath(router.config.BasePath, "/assets")
	data.CSRFToken = csrf.Token(r)

	err := ui.Render(w, data)
	if err != nil {
		return fmt.Errorf("error rendering UI template: %w", err)
	}

	return nil
}

func (router *router) ensureLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := session.Get[*auth.Account](r.Context(), "account")
		if !ok {
			router.redirectTo(w, r, "login")
			return
		}
		next.ServeHTTP(w, r)
	})
}
