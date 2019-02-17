package models

import (
	"github.com/jinzhu/gorm"
)

type Image struct {
	Name       string
	Type       string
	Path       string
	VoiceActor string
}

func FindImage(db *gorm.DB) []Image {
	var images []Image
	db.Table("image").Find(&images)
	return images
}
