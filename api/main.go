package api

import (
	"net/http"

	"github.com/chilipizdrick/muzek-server/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const API_ROUTE = "/api/v1"

func AssignRouteHandlers(r *gin.Engine, db *gorm.DB) {
	group := r.Group(API_ROUTE)
	groupAuth := r.Group(API_ROUTE, AuthMiddleware(db))

	assignAuthRouteHandlers(group, db)
	assignAlbumsRouteHandlers(groupAuth, db)
	assignArtistsRouteHandlers(groupAuth, db)
	assignPlaylistsRouteHandlers(groupAuth, db)
	assignSearchRouteHandlers(groupAuth, db)
	assignTracksRouteHandlers(groupAuth, db)
	assignUsersRouteHandlers(groupAuth, db)
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionId := session.Get("userId")
		if sessionId == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error: Error{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized.",
				},
			})
			return
		}
		userId := sessionId.(uint)
		_, err := database.GetUserByID(db, userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error: Error{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized.",
				},
			})
			return
		}
		c.Next()
	}
}
