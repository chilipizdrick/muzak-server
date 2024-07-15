package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const API_ROUTE = "/api/v1"

func AssignRouteHandlers(r *gin.Engine, db *gorm.DB) {
	group := r.Group(API_ROUTE)

	assignAlbumsRouteHandlers(group, db)
	assignArtistsRouteHandlers(group, db)
	assignPlaylistsRouteHandlers(group, db)
	assignSearchRouteHandlers(group, db)
	assignTracksRouteHandlers(group, db)
	assignUsersRouteHandlers(group, db)
}
