package database

import (
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username            string            `gorm:"type:string"`
	PasswordHash        string            `gorm:"type:string"`
	PlaylistIDs         pq.Int64Array     `gorm:"type:integer[]"`
	PlayerState         *PlayerStateModel `gorm:"embedded;embeddedPrefix:player_"`
	IsPlayerStatePublic bool              `gorm:"type:bool"`
	TSV                 string            `gorm:"->;type:tsvector GENERATED ALWAYS AS (to_tsvector('simple', username)) STORED;index:,type:GIN"`
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) AfterDelete(db *gorm.DB) (err error) {
	db.Where("owner_id = ?", m.ID).Delete(&PlaylistModel{})
	return
}

type PlayerStateModel struct {
	TrackID                 uint    `gorm:"type:uint"`
	IsPlaying               bool    `gorm:"type:bool"`
	Progress                uint    `gorm:"type:uint"` // In seconds
	Device                  string  `gorm:"type:string"`
	IsShuffleEnabled        bool    `gorm:"type:bool"`
	IsRepeatPlaylistEnabled bool    `gorm:"type:bool"`
	IsRepeatTrackEnabled    bool    `gorm:"type:bool"`
	Volume                  float64 `gorm:"type:float"` // From 0.0 to 1.0
}

type User struct {
	ID                  uint
	Username            string
	PasswordHash        string
	PlaylistIDs         []uint
	PlayerState         *PlayerState
	IsPlayerStatePublic bool
}

type UserExpanded struct {
	ID                  uint
	Username            string
	PasswordHash        string
	Playlists           []Playlist
	PlayerState         *PlayerState
	IsPlayerStatePublic bool
}

type PlayerState struct {
	Track                   Track
	IsPlaying               bool
	Progress                uint // In seconds
	Device                  string
	IsShuffleEnabled        bool
	IsRepeatPlaylistEnabled bool
	IsRepeatTrackEnabled    bool
	Volume                  float64 // From 0.0 to 1.0
}

func UserModelToUser(db *gorm.DB, userModel UserModel) (*User, error) {
	playerState, err := PlayerStateModelToPlayerState(db, userModel.PlayerState)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:                  userModel.ID,
		Username:            userModel.Username,
		PasswordHash:        userModel.PasswordHash,
		PlaylistIDs:         utils.PQInt64ArrayPtrToUIntSlice(userModel.PlaylistIDs),
		PlayerState:         playerState,
		IsPlayerStatePublic: userModel.IsPlayerStatePublic,
	}
	return &user, nil
}

func UserToUserExpanded(db *gorm.DB, user User) (*UserExpanded, error) {
	playlists, err := GetPlaylistsByIDs(db, user.PlaylistIDs)
	if err != nil {
		return nil, err
	}

	userExpanded := UserExpanded{
		ID:                  user.ID,
		Username:            user.Username,
		PasswordHash:        user.PasswordHash,
		Playlists:           playlists,
		PlayerState:         user.PlayerState,
		IsPlayerStatePublic: user.IsPlayerStatePublic,
	}
	return &userExpanded, nil
}

func UserModelToUserExpanded(db *gorm.DB, userModel UserModel) (*UserExpanded, error) {
	user, err := UserModelToUser(db, userModel)
	if err != nil {
		return nil, err
	}
	return UserToUserExpanded(db, *user)
}

func PlayerStateModelToPlayerState(db *gorm.DB, playerStateModel *PlayerStateModel) (*PlayerState, error) {
	if playerStateModel == nil {
		return nil, nil
	}

	track, err := GetTrackByID(db, playerStateModel.TrackID)
	if err != nil {
		return nil, err
	}

	playerState := PlayerState{
		Track:                   *track,
		IsPlaying:               playerStateModel.IsPlaying,
		Progress:                playerStateModel.Progress,
		Device:                  playerStateModel.Device,
		IsShuffleEnabled:        playerStateModel.IsShuffleEnabled,
		IsRepeatPlaylistEnabled: playerStateModel.IsRepeatPlaylistEnabled,
		IsRepeatTrackEnabled:    playerStateModel.IsRepeatTrackEnabled,
		Volume:                  playerStateModel.Volume,
	}

	return &playerState, nil
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

	user, err := UserModelToUser(db, *userModel)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserExpandedByID(db *gorm.DB, id uint) (*UserExpanded, error) {
	user, err := GetUserByID(db, id)
	if err != nil {
		return nil, err
	}

	userExpanded, err := UserToUserExpanded(db, *user)
	if err != nil {
		return nil, err
	}
	return userExpanded, nil
}

func GetUsersByIDs(db *gorm.DB, ids []uint) ([]User, error) {
	var userModels []UserModel
	if err := db.Where("id IN ?", ids).Find(&userModels).Error; err != nil {
		return nil, err
	}

	var users []User
	if userModels != nil {
		users = make([]User, len(userModels))
		for i, e := range userModels {
			user, err := UserModelToUser(db, e)
			if err != nil {
				return nil, err
			}
			users[i] = *user
		}
	}

	return users, nil
}
