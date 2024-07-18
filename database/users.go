package database

import (
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
	PublicPlaylists   bool              `gorm:"type:bool"`
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
	PublicPlaylists   bool
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
	PublicPlaylists   bool
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
