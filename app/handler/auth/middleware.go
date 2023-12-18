package auth

import (
	"context"
	"net/http"
	"strings"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
	"fmt"
)

var contextKey = new(struct{})

// Auth by header
func Middleware(app *app.App) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// ヘッダーから Username を取り出すだけの超安易な認証
			a := r.Header.Get("Authentication")

			fmt.Println("Authentication Header:", a) // ヘッダー内容を表示
			pair := strings.SplitN(a, " ", 2)
			if len(pair) < 2 {
				fmt.Println("Invalid Authentication Header:", a) // ヘッダーが不正な場合の表示
				httperror.Error(w, http.StatusUnauthorized)
				return
			}

			authType := pair[0]

			fmt.Println("Auth Type:", authType) // 認証タイプを表示
			if !strings.EqualFold(authType, "username") {
				fmt.Println("Invalid Auth Type:", authType) // 認証タイプが不正な場合の表示
				httperror.Error(w, http.StatusUnauthorized)
				return
			}

			username := pair[1]

			fmt.Println("Username:", username) // ユーザ名を表示
			if account, err := app.Dao.Account().FindByUsername(ctx, username); err != nil {
				fmt.Println("Error finding account:", err) // アカウント検索エラーを表示
				httperror.InternalServerError(w, err)
				return
			} else if account == nil {
				httperror.Error(w, http.StatusUnauthorized)
				return
			} else {
				fmt.Println("Account found:", account.Username) // アカウントが見つかった場合の表示
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey, account)))
			}
		})
	}
}

// Read Account data from authorized request
func AccountOf(r *http.Request) *object.Account {
	if cv := r.Context().Value(contextKey); cv == nil {
		fmt.Println("cv is nil")
		return nil

	} else if account, ok := cv.(*object.Account); !ok {
		fmt.Println("cv is not *object.Account")
		return nil

	} else {
		fmt.Println("cv is *object.Account")
		return account

	}
}
