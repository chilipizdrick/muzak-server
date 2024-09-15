package api

import (
	"log"
	"net/http"

	"github.com/chilipizdrick/muzek-server/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const SEARCH_API_ROUTE = "/search"

func assignSearchRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(SEARCH_API_ROUTE)

	group.GET("", searchWrapper(db))
}

type SearchRequest struct {
	Type  string `validate:"required,oneof=album artist playlist track user"`
	Query string `validate:"required,min=1"`
}

func searchWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		request := SearchRequest{
			Type:  c.Query("type"),
			Query: c.Query("q"),
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(request); err != nil {
			log.Printf("[TRACE] Failed to validate search request: %s", err)
			badRequestResponse(c, "Invalid search request.")
			return
		}

		switch request.Type {
		case "album":
			var albumModels []database.AlbumModel
			db.Limit(10).Where("tsv @@ to_tsquery('simple', ?)", request.Query).Find(&albumModels)

			albums := make([]AlbumExpanded, len(albumModels))
			for i, e := range albumModels {
				album, err := database.AlbumModelToAlbumExpanded(db, e)
				if err != nil {
					log.Printf("[ERROR] Failed to convert database.Album to database.AlbumExpanded: %s", err)
					internalServerErrorResponse(c, "Internal server error.")
					return
				}
				albums[i] = DBAlbumExpandedToAPIAlbumExpanded(*album)
			}
			c.JSON(http.StatusOK, albums)
			return

		case "artist":
			var artistModels []database.ArtistModel
			db.Limit(10).Where("tsv @@ to_tsquery('simple', ?)", request.Query).Find(&artistModels)

			artists := make([]ArtistExpanded, len(artistModels))
			for i, e := range artistModels {
				artist, err := database.ArtistModelToArtistExpanded(db, e)
				if err != nil {
					log.Printf("[ERROR] Failed to convert database.artist to database.artistExpanded: %s", err)
					internalServerErrorResponse(c, "Internal server error.")
					return
				}
				artists[i] = DBArtistExpandedToAPIArtistExpanded(*artist)
			}
			c.JSON(http.StatusOK, artists)
			return

		case "playlist":
			var playlistModels []database.PlaylistModel
			db.Limit(10).Where("tsv @@ to_tsquery('simple', ?)", request.Query).Find(&playlistModels)

			playlists := make([]PlaylistExpanded, len(playlistModels))
			for i, e := range playlistModels {
				playlist, err := database.PlaylistModelToPlaylistExpanded(db, e)
				if err != nil {
					log.Printf("[ERROR] Failed to convert database.playlist to database.playlistExpanded: %s", err)
					internalServerErrorResponse(c, "Internal server error.")
					return
				}
				playlists[i] = DBPlaylistExpandedToAPIPlaylistExpanded(*playlist)
			}
			c.JSON(http.StatusOK, playlists)
			return

		case "track":
			var trackModels []database.TrackModel
			db.Limit(10).Where("tsv @@ to_tsquery('simple', ?)", request.Query).Find(&trackModels)

			tracks := make([]TrackExpanded, len(trackModels))
			for i, e := range trackModels {
				track, err := database.TrackModelToTrackExpanded(db, e)
				if err != nil {
					log.Printf("[ERROR] Failed to convert database.track to database.trackExpanded: %s", err)
					internalServerErrorResponse(c, "Internal server error.")
					return
				}
				tracks[i] = DBTrackExpandedToAPITrackExpanded(*track)
			}
			c.JSON(http.StatusOK, tracks)
			return

		case "user":
			var userModels []database.UserModel
			db.Limit(10).Where("tsv @@ to_tsquery('simple', ?)", request.Query).Find(&userModels)

			users := make([]UserExpanded, len(userModels))
			for i, e := range userModels {
				user, err := database.UserModelToUserExpanded(db, e)
				if err != nil {
					log.Printf("[ERROR] Failed to convert database.user to database.userExpanded: %s", err)
					internalServerErrorResponse(c, "Internal server error.")
					return
				}
				users[i] = DBUserExpandedToAPIUserExpanded(*user)
			}
			c.JSON(http.StatusOK, users)
			return

		default:
			panic("unreachable")
		}
	}
}
