package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Title     string         `gorm:"type:string"`
	ArtistIDs *pq.Int64Array `gorm:"type:integer[]"`
	TrackIDs  *pq.Int64Array `gorm:"type:integer[]"`
}
