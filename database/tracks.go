package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Track struct {
	gorm.Model
	Title     string         `gorm:"type:string"`
	ArtistIDs *pq.Int64Array `gorm:"type:integer[]"`
	AlbumID   *uint          `gorm:"type:uint"`
	Genre     *string        `gorm:"type:string"`
	Duration  uint           `gorm:"type:uint"`
}

func GetTrackByIDFromDB(db *gorm.DB, id uint) (*Track, error) {
	var track Track
	err := db.First(&track, id).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}
