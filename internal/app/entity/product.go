package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID           uint32    `json:"id" bun:"id,notnull,unique,autoincrement"`
	GUID         uuid.UUID `json:"product_guid" bun:"type:uuid,notnull,pk"`
	Name         string    `json:"name" bun:"name,notnull,unique"`
	Description  string    `json:"description" bun:"description,default:null"`
	Price        float32   `json:"price" bun:"price,notnull"`
	CategoryGUID uuid.UUID `json:"category_guid" bun:"type:uuid,notnull"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt    time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
}

type RequestProductCreate struct {
	Name         string    `json:"name" binding:"required,min=2,max=255"`
	Price        float32   `json:"price" binding:"required,gt=0"`
	CategoryGUID uuid.UUID `json:"category_guid" binding:"required,uuid"`
	Description  string    `json:"description" binding:"omitempty,min=2,max=255"`
}

type RequestProductGetList struct {
	CategoryGUID uuid.UUID `json:"category_uuid" binding:"uuid,omitempty"`
	MinPrice     float32   `json:"min_price" binding:"required,ltfield=MaxPrice"`
	MaxPrice     float32   `json:"max_price" binding:"required,gtfield=MinPrice"`
}

type ResponseProductCreate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Price        float32   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  string    `json:"description"`
}

type RequestProductUpdate struct {
	Name         string    `json:"name" binding:"omitempty,min=2,max=255"`
	Price        float32   `json:"price" binding:"omitempty,gt=0"`
	CategoryGUID uuid.UUID `json:"category_guid" binding:"omitempty,uuid"`
	Description  string    `json:"description" binding:"omitempty,min=2,max=255"`
}
