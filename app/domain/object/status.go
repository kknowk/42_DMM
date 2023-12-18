package object

import (
	"time"
)

type (
	StatusID = int64

	Status struct {
		ID			StatusID	`db:"id"`
		AccountID	AccountID	`db:"account_id"`
		Content	string			`db:"content"`

		Username string `db:"username"`
		// The time the status was created
		CreateAt DateTime `db:"create_at"`

		// The time the status was updated
		UpdateAt DateTime `db:"update_at"`

		// The time the status was deleted
		DeleteAt DateTime `db:"delete_at"`
	}
)

func NewDateTime(t time.Time) DateTime {
	return DateTime{Time: t}
}

func (s *Status) GetID() int64 {
	return s.ID
}

func (s *Status) GetAccountID() int64 {
	return s.AccountID
}

func (s *Status) GetContent() string {
	return s.Content
}

func (s *Status) GetCreateAt() DateTime {
	return s.CreateAt
}