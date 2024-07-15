package database

import (
	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	Name     string
	OwnerID  uint
	IsPublic bool
}
