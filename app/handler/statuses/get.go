package statuses

import (
	"encoding/json"
	"net/http"
	"strconv"

	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)


// Handle request for `GET /v1/accounts/{username}`
func (h *handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	// パスからIDを取得
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// データベースからステータスを取得
	// fmt.Println("id: ", id)
	status, err := h.app.Dao.Status().FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	account, err := h.app.Dao.Account().FindByUsername(r.Context(), status.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// 応答用のStatusResponseを作成
	response := GetStatus{
		ID:				id,
		Account:		account,
		Content:		status.Content,
		CreateAt:		status.CreateAt,
		MediaAttachments:	[]MediaAttachment{}, // 空のメディア添付情報
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
