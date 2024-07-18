package database

import (
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TrackModel struct {
	gorm.Model
	Title     string        `gorm:"type:string"`
	ArtistIDs pq.Int64Array `gorm:"type:integer[]"`
	AlbumID   uint          `gorm:"type:uint"`
	Duration  uint          `gorm:"type:uint"`
}

func (TrackModel) TableName() string {
	return "tracks"
}

type Track struct {
	ID        uint
	Title     string
	ArtistIDs []uint
	AlbumID   uint
	Duration  uint
}

type TrackExpanded struct {
	ID       uint
	Title    string
	Artists  []Artist
	Album    Album
	Duration uint
}

func TrackModelToTrack(trackModel TrackModel) Track {
	return Track{
		ID:        trackModel.ID,
		Title:     trackModel.Title,
		ArtistIDs: utils.PQInt64ArrayPtrToUIntSlice(trackModel.ArtistIDs),
		AlbumID:   trackModel.AlbumID,
		Duration:  trackModel.Duration,
	}
}

func GetTrackModelByID(db *gorm.DB, id uint) (*TrackModel, error) {
	var trackModel TrackModel
	if err := db.First(&trackModel, id).Error; err != nil {
		return nil, err
	}
	return &trackModel, nil
}

func GetTrackByID(db *gorm.DB, id uint) (*Track, error) {
	trackModel, err := GetTrackModelByID(db, id)
	if err != nil {
		return nil, err
	}
	track := TrackModelToTrack(*trackModel)
	return &track, nil
}

func GetTrackExpandedByID(db *gorm.DB, id uint) (*TrackExpanded, error) {
	track, err := GetTrackByID(db, id)
	if err != nil {
		return nil, err
	}

	artists, err := GetArtistsByIDs(db, track.ArtistIDs)
	if err != nil {
		return nil, err
	}

	album, err := GetAlbumByID(db, track.AlbumID)
	if err != nil {
		return nil, err
	}

	trackExpanded := TrackExpanded{
		ID:       track.ID,
		Title:    track.Title,
		Artists:  artists,
		Album:    *album,
		Duration: track.Duration,
	}

	return &trackExpanded, nil
}

func GetTracksByIDs(db *gorm.DB, ids []uint) ([]Track, error) {
	var trackModels []TrackModel
	if err := db.Where(ids).Find(&trackModels).Error; err != nil {
		return nil, err
	}

	var tracks []Track
	if trackModels != nil {
		tracks = make([]Track, len(trackModels))
		for i, e := range trackModels {
			tracks[i] = TrackModelToTrack(e)
		}
	}

	return tracks, nil
}
