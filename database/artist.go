package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Artist struct {
	gorm.Model
	Name       string
	AlbumIDs   pq.Int64Array `gorm:"type:integer[]"`
	TrackIDs   pq.Int64Array `gorm:"type:integer[]"`
	IsVerified bool
}
