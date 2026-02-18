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
	Name         string    `json:"name"`
	Price        float32   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  string    `json:"description"`
}

type RequestProductGetList struct {
	CategoryGUID uuid.UUID `json:"category_uuid"`
	MinPrice     float32   `json:"min_price"`
	MaxPrice     float32   `json:"max_price"`
}

type ResponseProductCreate struct {
	GUID         uuid.UUID `json:"guid"`
	Name         string    `json:"name"`
	Price        float32   `json:"price"`
	CategoryGUID uuid.UUID `json:"category_guid"`
	Description  string    `json:"description"`
}
