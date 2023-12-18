package accounts

import (
	"net/http"

	"yatter-backend-go/app/app"

	"github.com/go-chi/chi"

	"yatter-backend-go/app/handler/accounts"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Post("/", h.Create)

	// TODO: Add other APIs
	r.Mount("/v1/accounts", accounts.NewRouter(app))

	return r
}
