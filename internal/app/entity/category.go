package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Category struct {
	bun.BaseModel `bun:"table:categories"`

	ID        uint32    `json:"id" bun:"id,notnull,unique,autoincrement"`
	GUID      uuid.UUID `json:"guid" bun:"type: uuid,pk,notnull"`
	Name      string    `json:"name" bun:"name,notnull,unique"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`
}

type RequestCategoryCreate struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}

type ResponseCategoryCreate struct {
	GUID      uuid.UUID `json:"guid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
