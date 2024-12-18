package model

import (
	"time"
)

type ImageTable struct {
	ID               uint      `gorm:"primaryKey"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ImageIdentifier  string    `json:"image_identifier" gorm:"uniqueIndex;size:20"`
	OriginalFileName string    `json:"original_file_name" gorm:"index"`
	ImageType        string    `json:"image_type"`
	ImageFileData    []byte    `json:"image_file_data"`
}
