package database

import (
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AlbumModel struct {
	gorm.Model
	Title     string        `gorm:"type:string"`
	ArtistIDs pq.Int64Array `gorm:"type:integer[]"`
	TrackIDs  pq.Int64Array `gorm:"type:integer[]"`
	TSV       string        `gorm:"->;type:tsvector GENERATED ALWAYS AS (to_tsvector('simple', title)) STORED;index:,type:GIN"`
}

func (AlbumModel) TableName() string {
	return "albums"
}

type Album struct {
	ID        uint
	Title     string
	ArtistIDs []uint
	TrackIDs  []uint
}

type AlbumExpanded struct {
	ID      uint
	Title   string
	Artists []Artist
	Tracks  []Track
}

func AlbumModelToAlbum(albumModel AlbumModel) Album {
	return Album{
		ID:        albumModel.ID,
		Title:     albumModel.Title,
		ArtistIDs: utils.PQInt64ArrayPtrToUIntSlice(albumModel.ArtistIDs),
		TrackIDs:  utils.PQInt64ArrayPtrToUIntSlice(albumModel.TrackIDs),
	}
}

func AlbumToAlbumExpanded(db *gorm.DB, album Album) (*AlbumExpanded, error) {
	artists, err := GetArtistsByIDs(db, album.ArtistIDs)
	if err != nil {
		return nil, err
	}

	tracks, err := GetTracksByIDs(db, album.TrackIDs)
	if err != nil {
		return nil, err
	}

	albumExpanded := AlbumExpanded{
		ID:      album.ID,
		Title:   album.Title,
		Artists: artists,
		Tracks:  tracks,
	}
	return &albumExpanded, nil
}

func AlbumModelToAlbumExpanded(db *gorm.DB, albumModel AlbumModel) (*AlbumExpanded, error) {
	return AlbumToAlbumExpanded(db, AlbumModelToAlbum(albumModel))
}

func GetAlbumModelByID(db *gorm.DB, id uint) (*AlbumModel, error) {
	var albumModel AlbumModel
	if err := db.First(&albumModel, id).Error; err != nil {
		return nil, err
	}
	return &albumModel, nil
}

func GetAlbumByID(db *gorm.DB, id uint) (*Album, error) {
	albumModel, err := GetAlbumModelByID(db, id)
	if err != nil {
		return nil, err
	}

	album := AlbumModelToAlbum(*albumModel)
	return &album, nil
}

func GetAlbumExpandedByID(db *gorm.DB, id uint) (*AlbumExpanded, error) {
	album, err := GetAlbumByID(db, id)
	if err != nil {
		return nil, err
	}

	albumExpanded, err := AlbumToAlbumExpanded(db, *album)
	if err != nil {
		return nil, err
	}
	return albumExpanded, nil
}

func GetAlbumsByIDs(db *gorm.DB, ids []uint) ([]Album, error) {
	var albumModels []AlbumModel
	if err := db.Where("id IN ?", ids).Find(&albumModels).Error; err != nil {
		return nil, err
	}

	var albums []Album
	if albumModels != nil {
		albums = make([]Album, len(albumModels))
		for i, e := range albumModels {
			albums[i] = AlbumModelToAlbum(e)
		}
	}
	return albums, nil
}
