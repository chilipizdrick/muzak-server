package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string        `gorm:"type:string"`
	PasswordHash      string        `gorm:"type:string"`
	PlaylistIDs       pq.Int64Array `gorm:"type:integer[]"`
	PlayerState       PlayerState   `gorm:"embedded;embeddedPrefix:player_"`
	PublicPlayerState bool          `gorm:"type:bool"`
}

type PlayerState struct {
	TrackID               *uint    `gorm:"type:uint"`
	Progress              *uint    `gorm:"type:uint"` // In seconds
	Device                *string  `gorm:"type:string"`
	ShuffleEnabled        *bool    `gorm:"type:bool"`
	RepeatPlaylistEnabled *bool    `gorm:"type:bool"`
	RepeatTrackEnabled    *bool    `gorm:"type:bool"`
	Volume                *float64 `gorm:"type:float"` // From 0.0 to 1.0
}
