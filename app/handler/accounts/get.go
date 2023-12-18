package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
	"github.com/go-chi/chi"
)

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	// ユーザー名をログに出力
	account, err := h.app.Dao.Account().FindByUsername(r.Context(), username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if account == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}