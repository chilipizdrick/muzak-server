package database

import "gorm.io/gorm"

func AutoMigrateSchemas(db *gorm.DB) {
    db.AutoMigrate(&Album{})
    db.AutoMigrate(&Artist{})
    db.AutoMigrate(&Playlist{})
    db.AutoMigrate(&Track{})
    db.AutoMigrate(&User{})
}
