package database

import "gorm.io/gorm"

func AutoMigrateSchemas(db *gorm.DB) {
	db.AutoMigrate(&AlbumModel{})
	db.AutoMigrate(&ArtistModel{})
	db.AutoMigrate(&PlaylistModel{})
	db.AutoMigrate(&TrackModel{})
	db.AutoMigrate(&UserModel{})
}
