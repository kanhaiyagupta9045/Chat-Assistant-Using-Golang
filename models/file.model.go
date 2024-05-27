package models

import (
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
