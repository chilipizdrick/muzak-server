package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AlbumModel struct {
	gorm.Model
	Title     string         `gorm:"type:string"`
	ArtistIDs *pq.Int64Array `gorm:"type:integer[]"`
	TrackIDs  *pq.Int64Array `gorm:"type:integer[]"`
}

type Album struct {
	ID    uint
	Title string
	ArtistIDs []uint
	TrackIDs  []uint
}

type AlbumExpanded struct {
	ID      uint
	Title   string
	Artists []Artist
	Tracks  []Track
}
