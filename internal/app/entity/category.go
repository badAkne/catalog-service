package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories"`

	ID        uint32    `json:"id" bun:"id,notnull,unique"`
	GUID      string    `json:"guid" bun:"guid,pk,notnull"`
	Name      string    `json:"name" bun:"name,notnull,unique"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`
}

type CategoryResponse struct {
}
