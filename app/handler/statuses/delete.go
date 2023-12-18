package statuses

import (
	"net/http"
	"strconv"
	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// DeleteStatus は指定されたIDのステータスを削除するためのリクエストを処理します
func (h *handler) DeleteStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	status, err := h.app.Dao.Status().FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	err = h.app.Dao.Status().DeleteStatus(r.Context(), status)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}