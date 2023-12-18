package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

func (r *account) Add(ctx context.Context, acc *object.Account) error {
	// SQLクエリの準備
	query := `INSERT INTO account (username, password_hash, display_name, avatar, header, note, create_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	// SQLクエリの実行
	result, err := r.db.ExecContext(ctx, query, acc.Username, acc.PasswordHash, acc.DisplayName, acc.Avatar, acc.Header, acc.Note, acc.CreateAt)
	if err != nil {
		return fmt.Errorf("failed to add account: %w", err)
	}
	
	// 新しく追加されたアカウントのIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	// IDをアカウントオブジェクトに設定
	acc.ID = object.AccountID(id)
	return nil
}


// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := new(object.Account)
	err := r.db.QueryRowxContext(ctx, "select * from account where username = ?", username).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}
