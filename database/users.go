package database

import (
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username          string            `gorm:"type:string"`
	PasswordHash      string            `gorm:"type:string"`
	PlaylistIDs       pq.Int64Array     `gorm:"type:integer[]"`
	PlayerState       *PlayerStateModel `gorm:"embedded;embeddedPrefix:player_"`
	PublicPlayerState bool              `gorm:"type:bool"`
}

func (UserModel) TableName() string {
	return "users"
}

type PlayerStateModel struct {
	TrackID               uint    `gorm:"type:uint"`
	Progress              uint    `gorm:"type:uint"` // In seconds
	Device                string  `gorm:"type:string"`
	ShuffleEnabled        bool    `gorm:"type:bool"`
	RepeatPlaylistEnabled bool    `gorm:"type:bool"`
	RepeatTrackEnabled    bool    `gorm:"type:bool"`
	Volume                float64 `gorm:"type:float"` // From 0.0 to 1.0
}

type User struct {
	ID                uint
	Username          string
	PasswordHash      string
	PlaylistIDs       []uint
	PlayerState       *PlayerState
	PublicPlayerState bool
}

type PlayerState struct {
	TrackID               uint
	Progress              uint // In seconds
	Device                string
	ShuffleEnabled        bool
	RepeatPlaylistEnabled bool
	RepeatTrackEnabled    bool
	Volume                float64 // From 0.0 to 1.0
}

type UserExpanded struct {
	ID                uint
	Username          string
	PasswordHash      string
	Playlists         []Playlist
	PlayerState       *PlayerStateExpanded
	PublicPlayerState bool
}

type PlayerStateExpanded struct {
	Track                 Track
	Progress              uint // In seconds
	Device                string
	ShuffleEnabled        bool
	RepeatPlaylistEnabled bool
	RepeatTrackEnabled    bool
	Volume                float64 // From 0.0 to 1.0
}

func UserModelToUser(userModel UserModel) User {
	return User{
		ID:                userModel.ID,
		Username:          userModel.Username,
		PasswordHash:      userModel.PasswordHash,
		PlaylistIDs:       utils.PQInt64ArrayPtrToUIntSlice(userModel.PlaylistIDs),
		PlayerState:       PlayerStateModelToPlayerState(userModel.PlayerState),
		PublicPlayerState: userModel.PublicPlayerState,
	}
}

func PlayerStateModelToPlayerState(playerStateModel *PlayerStateModel) *PlayerState {
	if playerStateModel == nil {
		return nil
	}

	playerState := PlayerState{
		TrackID:               playerStateModel.TrackID,
		Progress:              playerStateModel.Progress,
		Device:                playerStateModel.Device,
		ShuffleEnabled:        playerStateModel.ShuffleEnabled,
		RepeatPlaylistEnabled: playerStateModel.RepeatPlaylistEnabled,
		RepeatTrackEnabled:    playerStateModel.RepeatTrackEnabled,
		Volume:                playerStateModel.Volume,
	}

	return &playerState
}

func PlayerStateToPlayerStateExpanded(db *gorm.DB, playerState *PlayerState) (*PlayerStateExpanded, error) {
	if playerState == nil {
		return nil, nil
	}

	track, err := GetTrackByID(db, playerState.TrackID)
	if err != nil {
		return nil, err
	}

	playerStateExpanded := PlayerStateExpanded{
		Track:                 *track,
		Progress:              playerState.Progress,
		Device:                playerState.Device,
		ShuffleEnabled:        playerState.ShuffleEnabled,
		RepeatPlaylistEnabled: playerState.RepeatPlaylistEnabled,
		RepeatTrackEnabled:    playerState.RepeatTrackEnabled,
		Volume:                playerState.Volume,
	}

	return &playerStateExpanded, nil
}

func GetUserModelByID(db *gorm.DB, id uint) (*UserModel, error) {
	var userModel UserModel
	if err := db.First(&userModel, id).Error; err != nil {
		return nil, err
	}
	return &userModel, nil
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	userModel, err := GetUserModelByID(db, id)
	if err != nil {
		return nil, err
	}
	user := UserModelToUser(*userModel)
	return &user, nil
}

func GetUserExpandedByID(db *gorm.DB, id uint) (*UserExpanded, error) {
	user, err := GetUserByID(db, id)
	if err != nil {
		return nil, err
	}

	playerStateExpanded, err := PlayerStateToPlayerStateExpanded(db, user.PlayerState)
	if err != nil {
		return nil, err
	}

	userExpanded := UserExpanded{
		ID:                user.ID,
		Username:          user.Username,
		PasswordHash:      user.PasswordHash,
		PlayerState:       playerStateExpanded,
		PublicPlayerState: user.PublicPlayerState,
	}
	return &userExpanded, nil
}

func GetUsersByIDs(db *gorm.DB, ids []uint) ([]User, error) {
	var userModels []UserModel
	if err := db.Where(ids).Find(&userModels).Error; err != nil {
		return nil, err
	}

	var users []User
	if userModels != nil {
		users = make([]User, len(userModels))
		for i, e := range userModels {
			users[i] = UserModelToUser(e)
		}
	}

	return users, nil
}
