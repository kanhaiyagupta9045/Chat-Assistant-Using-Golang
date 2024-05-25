package models

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Filename string
	Content  []byte
}

type Assistant struct {
	AssistantName         string `json:"assistant_name" validate:"required"`
	AssistantInstructions string `json:"assistant_instruction" validate:"required"`
	Content               string `json:"content" validate:"required"`
}

type FileUploadRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Assistant Assistant             `form:"assistant" binding:"required"`
}
