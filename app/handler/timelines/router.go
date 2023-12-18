package timelines

import (
	"net/http"
	"yatter-backend-go/app/app"
	"github.com/go-chi/chi"

)
// Implementation of handler
type handler struct {
	app *app.App
}

func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()
	h := &handler{app: app}

	r.Get("/public", h.GetPublicTimeline)

	return r
}
