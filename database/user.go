package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username                    string
	PasswordHash                string
	PlaylistIDs                 pq.Int64Array `gorm:"type:integer[]"`
	PlayerTrackID               *uint
	PlayerProgress              *uint // In seconds
	PlayerDevice                *string
	PlayerShuffleEnabled        *bool
	PlayerRepeatPlaylistEnabled *bool
	PlayerRepeatTrackEnabled    *bool
	PlayerVolume                float64 // From 0.0 to 1.0
}
