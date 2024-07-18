package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ArtistModel struct {
	gorm.Model
	Name       string         `gorm:"type:string"`
	AlbumIDs   *pq.Int64Array `gorm:"type:integer[]"`
	TrackIDs   *pq.Int64Array `gorm:"type:integer[]"`
	IsVerified bool           `gorm:"type:bool"`
}

type Artist struct {
	ID         uint
	Name       string
	AlbumIDs     []uint
	TrackIDs     []uint
	IsVerified bool
}

type ArtistExpanded struct {
	ID         uint
	Name       string
	Albums     []Album
	Tracks     []Track
	IsVerified bool
}
