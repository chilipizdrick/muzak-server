package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PlaylistModel struct {
	gorm.Model
	Title    string         `gorm:"type:string"`
	OwnerID  uint           `gorm:"type:uint"`
	IsPublic bool           `gorm:"type:bool"`
	TrackIDs *pq.Int64Array `gorm:"type:integer[]"`
}

type Playlist struct {
	ID       uint
	Title    string
	OwnerID  uint
	IsPublic bool
	TrackIDs []uint
}

type PlaylistExpanded struct {
	ID       uint
	Title    string
	Owner    User
	IsPublic bool
	Tracks   []Track
}
