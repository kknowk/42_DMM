package statuses

import (
	"net/http"

	"yatter-backend-go/app/app"

	"github.com/go-chi/chi"
	"yatter-backend-go/app/handler/auth"
	// "yatter-backend-go/app/handler/statuses"

)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}

		r.Route("/", func(r chi.Router) {
		// 認証が不要なルート
		r.Get("/{id}", h.GetStatus)

		r.Delete("/{id}", h.DeleteStatus)

		// 認証が必要なルート
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware(app)) // このグループ内のルートに認証ミドルウェアを適用
			r.Post("/", h.CreateStatus) // このルートには認証を適用
		})
	})


	return r
}
