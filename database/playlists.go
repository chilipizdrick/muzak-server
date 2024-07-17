package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	Name     string         `gorm:"type:string"`
	OwnerID  uint           `gorm:"type:uint"`
	IsPublic bool           `gorm:"type:bool"`
	TrackIDs *pq.Int64Array `gorm:"type:integer[]"`
}
