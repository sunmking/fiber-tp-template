package model

import (
	"time"

	"fiber-blog/config"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_time"` // CreatedAt gorm will generate the value when insert
	UpdatedAt time.Time `json:"updated_time"` // UpdatedAt gorm will generate the value when insert
}

func db() *gorm.DB {
	return config.DB
}
