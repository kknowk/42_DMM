package accounts

import (
	"encoding/json"
	"net/http"
	"time"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/accounts`
type AddRequest struct {
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	DisplayName  *string `json:"display_name,omitempty"`
	Avatar       *string `json:"avatar,omitempty"`
	Header       *string `json:"header,omitempty"`
	Note         *string `json:"note,omitempty"`
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := &object.Account{
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Avatar:      req.Avatar,
		Header:      req.Header,
		Note:        req.Note,
		CreateAt:    object.DateTime{Time: time.Now()},
	}
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// ユーザー名が既に使用されているか確認
	existingAccount, err := h.app.Dao.Account().FindByUsername(r.Context(), req.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if existingAccount != nil {
		// 既に存在するユーザー名の場合はエラーを返す
		http.Error(w, "Username already in use", http.StatusConflict)
		return
	}

	// アカウントをデータベースに追加
	if err := h.app.Dao.Account().Add(r.Context(), account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	// 成功レスポンスを返す
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
