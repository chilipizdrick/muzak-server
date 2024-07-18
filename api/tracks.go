package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/chilipizdrick/muzek-server/database"
	"github.com/chilipizdrick/muzek-server/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const TRACKS_API_ROUTE = "/tracks"

func assignTracksRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(TRACKS_API_ROUTE)

	group.GET("/:id", getTrackByIDWrapper(db))
}

type Track struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	ArtistIDs *[]uint `json:"artistIds"`
	AlbumID   *uint   `json:"albumId"`
	Genre     *string `json:"genre"`
	Duration  uint    `json:"duration"`
	Href      string  `json:"href"`
}

func DBTrackToAPITrack(track database.TrackModel) Track {
	return Track{
		ID:        track.ID,
		Title:     track.Title,
		ArtistIDs: utils.PQInt64ArrayPtrToUIntSlice(track.ArtistIDs),
		AlbumID:   track.AlbumID,
		Genre:     track.Genre,
		Duration:  track.Duration,
		Href:      fmt.Sprintf("%s/tracks/%d/audio.ogg", os.Getenv("ASSETS_SERVER_URI"), track.ID),
	}
}

func getTrackByIDWrapper(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id64, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Printf("[INFO] Non integer id \"%v\" has been provided: %s", idString, err)
			c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{
				Error: Error{
					Status:  http.StatusInternalServerError,
					Message: "Invalid track id.",
				},
			})
			return
		}
		id := uint(id64)
		track, err := database.GetTrackByIDFromDB(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.IndentedJSON(http.StatusNotFound, ErrorResponse{
					Error: Error{
						Status:  http.StatusNotFound,
						Message: fmt.Sprintf("Track with id \"%d\" was not found.", id),
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
		c.IndentedJSON(http.StatusOK, DBTrackToAPITrack(*track))
	}
}
