package database

import (
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TrackModel struct {
	gorm.Model
	Title     string         `gorm:"type:string"`
	ArtistIDs *pq.Int64Array `gorm:"type:integer[]"`
	AlbumID   *uint          `gorm:"type:uint"`
	Genre     *string        `gorm:"type:string"`
	Duration  uint           `gorm:"type:uint"`
}

type Track struct {
	ID        uint
	Title     string
	ArtistIDs []uint
	AlbumID   *uint
	Genre     *string
	Duration  uint
}

type TrackExpanded struct {
	ID       uint
	Title    string
	Artists  []Artist
	Album    *Album
	Genre    *string
	Duration uint
}

func GetTrackModelByIDFromDB(db *gorm.DB, id uint) (*TrackModel, error) {
	var trackModel TrackModel
	if err := db.First(&trackModel, id).Error; err != nil {
		return nil, err
	}
	return &trackModel, nil
}

func GetTrackByIDFromDB(db *gorm.DB, id uint) (*Track, error) {
	trackModel, err := GetTrackModelByIDFromDB(db, id)
	if err != nil {
		return nil, err
	}
	track := Track{
		ID:        trackModel.ID,
		Title:     trackModel.Title,
		ArtistIDs: *utils.PQInt64ArrayPtrToUIntSlice(trackModel.ArtistIDs),
		AlbumID:   trackModel.AlbumID,
		Genre:     trackModel.Genre,
		Duration:  trackModel.Duration,
	}
	return &track, nil
}
