package statuses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

type StatusRequest struct {
	Status		string `json:"status"`
	MediaIDs	[]int  `json:"media_ids"`
}

type StatusResponse struct {
	ID					object.StatusID		`json:"id"`
	Account				*object.Account		`json:"account"`
	Content				string				`json:"content"`
	CreateAt			object.DateTime		`json:"create_at"`
	MediaAttachments	[]MediaAttachment	`json:"media_attachments"`
}

type MediaAttachment struct {
	ID			int64  `json:"id"`
	Type		string `json:"type"`
	URL			string `json:"url"`
	Description	string `json:"description"`
}

type GetStatus struct {
	ID					object.StatusID		`json:"id"`
	Account				*object.Account		`json:"account"`
	Content				string				`json:"content"`
	CreateAt			object.DateTime		`json:"create_at"`
	MediaAttachments	[]MediaAttachment	`json:"media_attachments"`
}

// CreateStatus は新しいステータスを作成するための POST リクエストを処理します
func (h *handler) CreateStatus(w http.ResponseWriter, r *http.Request) {
	// 認証チェック
	account := auth.AccountOf(r)
	if account == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// バリデーション: テキストが空でないことを確認
	fmt.Println("req.Content:", req.Status)
	if req.Status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	// バリデーション: テキストの長さを確認（例: 最大280文字）
	if len(req.Status) > 280 {
		http.Error(w, "Status is too long", http.StatusBadRequest)
		return
	}

	newStatus := object.Status{
		AccountID:	account.ID,
		Content:	req.Status,
		CreateAt:	object.NewDateTime(time.Now()),
		Username:	account.Username,
	}

	// データベースにステータスを保存
	if err := h.app.Dao.Status().AddStatus(r.Context(), &newStatus); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// 応答用のStatusResponseを作成
	response := StatusResponse{
			ID:				newStatus.ID,
			Account:		account,
			Content:		newStatus.Content,
			CreateAt:		newStatus.CreateAt,
			MediaAttachments:	[]MediaAttachment{}, // 空のメディア添付情報
		}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
