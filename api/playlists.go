package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chilipizdrick/muzek-server/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const PLAYLISTS_API_ROUTE = "/playlists"

func assignPlaylistsRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(PLAYLISTS_API_ROUTE)

	group.GET("/:id", getPlaylistByIDWrapper(db))
}

type Playlist struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	IsPublic bool   `json:"isPublic"`
}

type PlaylistExpanded struct {
	ID       uint    `json:"id"`
	Title    string  `json:"title"`
	Owner    User    `json:"owner"`
	IsPublic bool    `json:"isPublic"`
	Tracks   []Track `json:"tracks"`
}

func DBPlaylistToAPIPlaylist(playlist database.Playlist) Playlist {
	return Playlist{
		ID:       playlist.ID,
		Title:    playlist.Title,
		IsPublic: playlist.IsPublic,
	}
}

func DBPlaylistExpandedToAPIPlaylistExpanded(playlist database.PlaylistExpanded) PlaylistExpanded {
	tracks := make([]Track, len(playlist.Tracks))
	for i, e := range playlist.Tracks {
		tracks[i] = DBTrackToAPITrack(e)
	}

	return PlaylistExpanded{
		ID:       playlist.ID,
		Title:    playlist.Title,
		Owner:    DBUserToAPIUser(*playlist.Owner),
		IsPublic: playlist.IsPublic,
		Tracks:   tracks,
	}
}

func getPlaylistByIDWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id64, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Printf("[INFO] Non integer id \"%s\" has been provided: %s", idString, err)
			c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{
				Error: Error{
					Status:  http.StatusInternalServerError,
					Message: "Invalid artist id.",
				},
			})
			return
		}
		id := uint(id64)
		playlist, err := database.GetPlaylistExpandedByID(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.IndentedJSON(http.StatusNotFound, ErrorResponse{
					Error: Error{
						Status:  http.StatusNotFound,
						Message: fmt.Sprintf("Artist with id \"%d\" was not found.", id),
					},
				})
				return
			}
			c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{
				Error: Error{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error.",
				},
			})
			return
		}
		c.IndentedJSON(http.StatusOK, DBPlaylistExpandedToAPIPlaylistExpanded(*playlist))
	}
}
