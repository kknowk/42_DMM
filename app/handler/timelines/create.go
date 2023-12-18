package timelines

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/handler/httperror"

	"strconv"
	// "yatter-backend-go/app/domain/object"
)

const DefaultLimit = 40 // デフォルトの制限値
const MaxLimit = 80     // 最大許容制限値

func (h *handler) GetPublicTimeline(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータの取得
	onlyMediaStr := r.URL.Query().Get("only_media")
	maxIDStr := r.URL.Query().Get("max_id")
	sinceIDStr := r.URL.Query().Get("since_id")
	limitStr := r.URL.Query().Get("limit")

	var onlyMedia bool
	var maxID, sinceID int64
	var limit int
	var err error

	// onlyMedia の変換。空でなければ変換する
	if onlyMediaStr != "" {
		onlyMedia, err = strconv.ParseBool(onlyMediaStr)
		if err != nil {
			// 不正なクエリパラメータの場合のエラーハンドリング
			http.Error(w, "Invalid only_media parameter", http.StatusBadRequest)
			return
		}
	}

	// maxID の変換。空でなければ変換する
	if maxIDStr != "" {
		maxID, err = strconv.ParseInt(maxIDStr, 10, 64)
		if err != nil || maxID < 0 {
			http.Error(w, "Invalid max_id parameter", http.StatusBadRequest)
			return
		}
	}

	// sinceID の変換。空でなければ変換する
	if sinceIDStr != "" {
		sinceID, err = strconv.ParseInt(sinceIDStr, 10, 64)
		if err != nil || maxID < 0 {
			http.Error(w, "Invalid since_id parameter", http.StatusBadRequest)
			return
		}
	}

	// limit の変換。空でなければ変換する
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > MaxLimit {
			limit = DefaultLimit
		}
	} else {
		limit = DefaultLimit
	}

	statuses, err := h.app.Dao.Status().FindAllPublicStatuses(r.Context(), onlyMedia, maxID, sinceID, limit)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
