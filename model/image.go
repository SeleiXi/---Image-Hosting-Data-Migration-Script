package model

import (
	"time"
)

type NewImageTable struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ImageIdentifier  string    `json:"image_identifier" gorm:"uniqueIndex;size:20"`
	OriginalFileName string    `json:"original_file_name" gorm:"index"`
	ImageType        string    `json:"image_type"`
	ImageFileData    []byte    `json:"image_file_data"`
}

func (NewImageTable) TableName() string {
	return "image_table" // 映射新图床表名
}

type OriginalImageTable struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	OriginName string
	Path       string
	Name       string // 包括identifier和extension
}

func (OriginalImageTable) TableName() string {
	return "images" // 测试服的旧图床表名
}
