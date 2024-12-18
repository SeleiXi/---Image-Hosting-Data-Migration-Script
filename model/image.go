package model

import (
	"time"
)

type NewImageTable struct {
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ImageIdentifier  string    `json:"image_identifier" gorm:"uniqueIndex;size:20"`
	OriginalFileName string    `json:"original_file_name" gorm:"index"`
	ImageType        string    `json:"image_type"`
	ImageFileData    []byte    `json:"image_file_data"`
}

type OriginalImageTable struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	OriginName string
	Name       string // 包括identifier和extension
}
