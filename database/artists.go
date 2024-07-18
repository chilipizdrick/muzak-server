package database

import (
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ArtistModel struct {
	gorm.Model
	Name       string        `gorm:"type:string"`
	AlbumIDs   pq.Int64Array `gorm:"type:integer[]"`
	TrackIDs   pq.Int64Array `gorm:"type:integer[]"`
	IsVerified bool          `gorm:"type:bool"`
}

func (ArtistModel) TableName() string {
	return "artists"
}

type Artist struct {
	ID         uint
	Name       string
	AlbumIDs   []uint
	TrackIDs   []uint
	IsVerified bool
}

type ArtistExpanded struct {
	ID         uint
	Name       string
	Albums     []Album
	Tracks     []Track
	IsVerified bool
}

func ArtistModelToArtist(artistModel ArtistModel) Artist {
	return Artist{
		ID:         artistModel.ID,
		Name:       artistModel.Name,
		AlbumIDs:   utils.PQInt64ArrayPtrToUIntSlice(artistModel.AlbumIDs),
		TrackIDs:   utils.PQInt64ArrayPtrToUIntSlice(artistModel.TrackIDs),
		IsVerified: artistModel.IsVerified,
	}
}

func GetArtistModelByID(db *gorm.DB, id uint) (*ArtistModel, error) {
	var artistModel ArtistModel
	if err := db.First(&artistModel, id).Error; err != nil {
		return nil, err
	}
	return &artistModel, nil
}

func GetArtistByID(db *gorm.DB, id uint) (*Artist, error) {
	artistModel, err := GetArtistModelByID(db, id)
	if err != nil {
		return nil, err
	}
	artist := ArtistModelToArtist(*artistModel)
	return &artist, nil
}

func GetArtistExpandedByID(db *gorm.DB, id uint) (*ArtistExpanded, error) {
	artist, err := GetArtistByID(db, id)
	if err != nil {
		return nil, err
	}

	albums, err := GetAlbumsByIDs(db, artist.AlbumIDs)
	if err != nil {
		return nil, err
	}

	tracks, err := GetTracksByIDs(db, artist.AlbumIDs)
	if err != nil {
		return nil, err
	}

	artistExpanded := ArtistExpanded{
		ID:         artist.ID,
		Name:       artist.Name,
		Albums:     albums,
		Tracks:     tracks,
		IsVerified: artist.IsVerified,
	}
	return &artistExpanded, nil
}

func GetArtistsByIDs(db *gorm.DB, ids []uint) ([]Artist, error) {
	var artistModels []ArtistModel
	if err := db.Where(ids).Find(&artistModels).Error; err != nil {
		return nil, err
	}

	var artists []Artist
	if artistModels != nil {
		artists = make([]Artist, len(artistModels))
		for i, e := range artistModels {
			artists[i] = ArtistModelToArtist(e)
		}
	}

	return artists, nil
}
