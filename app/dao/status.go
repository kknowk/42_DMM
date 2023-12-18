package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type statusDao struct {
	db *sqlx.DB
}

func NewStatus(db *sqlx.DB) repository.Status {
	return &statusDao{db: db}
}

func (r *statusDao) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	// var status object.Status

	fmt.Println("id: ", id)
	entity := new(object.Status)
	err := r.db.QueryRowxContext(ctx, "SELECT * FROM status WHERE id = ?", id).StructScan(entity)
	fmt.Println("entity: ", entity)
	if err != nil {
		if err == sql.ErrNoRows {
			// ステータスが存在しない場合の処理
			return nil, fmt.Errorf("status not found")
		}
		return nil, err
	}
	return entity, nil
}

func (dao *statusDao) AddStatus(ctx context.Context, status *object.Status) error {
	// SQLクエリの準備と実行
	query := `INSERT INTO status (account_id, content, create_at, username) VALUES (?, ?, ?, ?)`

	// SQLクエリの実行
	result, err := dao.db.ExecContext(ctx, query, status.AccountID, status.Content, status.CreateAt, status.Username)
	if err != nil {
		return fmt.Errorf("failed to add status: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	status.ID = object.StatusID(id)
	return nil
}

func (dao *statusDao) DeleteStatus(ctx context.Context, status *object.Status) error {
	// SQLクエリの準備と実行
	query := `DELETE FROM status WHERE id = ?`

	// SQLクエリの実行
	result, err := dao.db.ExecContext(ctx, query, status.ID)
	if err != nil {
		return fmt.Errorf("failed to delete status: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	status.ID = object.StatusID(id)
	return nil
}

// FindAllPublicStatuses は全ての公開ステータスを取得します
func (dao *statusDao) FindAllPublicStatuses(ctx context.Context, onlyMedia bool, maxID int64, sinceID int64, limit int) ([]*object.Status, error) {
	if onlyMedia {
		// queryBuilder.WriteString(" WHERE media_attachments IS NOT NULL")
		return nil, nil
	}
	var statuses []*object.Status

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT * FROM status")

	var args []interface{}

	if maxID > 0 {
		if onlyMedia {
			queryBuilder.WriteString(" AND")
		} else {
			queryBuilder.WriteString(" WHERE")
		}
		queryBuilder.WriteString(" id < ?")
		args = append(args, maxID)
	}

	if sinceID > 0 {
		if onlyMedia || maxID > 0 {
			queryBuilder.WriteString(" AND")
		} else {
			queryBuilder.WriteString(" WHERE")
		}
		queryBuilder.WriteString(" id > ?")
		args = append(args, sinceID)
	}

	queryBuilder.WriteString(" ORDER BY create_at DESC")

	if limit > 0 {
		queryBuilder.WriteString(" LIMIT ?")
		args = append(args, limit)
	}

	query := queryBuilder.String()
	err := dao.db.SelectContext(ctx, &statuses, query, args...)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
