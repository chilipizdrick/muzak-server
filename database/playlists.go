package database

import (
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PlaylistModel struct {
	gorm.Model
	Title    string        `gorm:"type:string"`
	OwnerID  uint          `gorm:"type:uint"`
	IsPublic bool          `gorm:"type:bool"`
	TrackIDs pq.Int64Array `gorm:"type:integer[]"`
}

func (PlaylistModel) TableName() string {
	return "playlists"
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

func PlaylistModelToPlaylist(playlistModel PlaylistModel) Playlist {
	return Playlist{
		ID:       playlistModel.ID,
		Title:    playlistModel.Title,
		OwnerID:  playlistModel.OwnerID,
		IsPublic: playlistModel.IsPublic,
		TrackIDs: utils.PQInt64ArrayPtrToUIntSlice(playlistModel.TrackIDs),
	}
}

func GetPlaylistModelByID(db *gorm.DB, id uint) (*PlaylistModel, error) {
	var playlistModel PlaylistModel
	if err := db.First(&playlistModel, id).Error; err != nil {
		return nil, err
	}
	return &playlistModel, nil
}

func GetPlaylistByID(db *gorm.DB, id uint) (*Playlist, error) {
	playlistModel, err := GetPlaylistModelByID(db, id)
	if err != nil {
		return nil, err
	}
	playlist := PlaylistModelToPlaylist(*playlistModel)
	return &playlist, nil
}

func GetPlaylistExpandedByID(db *gorm.DB, id uint) (*PlaylistExpanded, error) {
	playlist, err := GetPlaylistByID(db, id)
	if err != nil {
		return nil, err
	}

	owner, err := GetUserByID(db, playlist.OwnerID)
	if err != nil {
		return nil, err
	}

	tracks, err := GetTracksByIDs(db, playlist.TrackIDs)
	if err != nil {
		return nil, err
	}

	playlistExpanded := PlaylistExpanded{
		ID:       playlist.ID,
		Title:    playlist.Title,
		Owner:    *owner,
		IsPublic: playlist.IsPublic,
		Tracks:   tracks,
	}
	return &playlistExpanded, nil
}

func GetPlaylistsByIDs(db *gorm.DB, ids []uint) ([]Playlist, error) {
	var playlistModels []PlaylistModel
	if err := db.Where(ids).Find(&playlistModels).Error; err != nil {
		return nil, err
	}

	var playlists []Playlist
	if playlistModels != nil {
		playlists = make([]Playlist, len(playlistModels))
		for i, e := range playlistModels {
			playlists[i] = PlaylistModelToPlaylist(e)
		}
	}

	return playlists, nil
}
