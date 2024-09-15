package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/chilipizdrick/muzek-server/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const TRACKS_API_ROUTE = "/tracks"

func assignTracksRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(TRACKS_API_ROUTE)

	group.GET("/:id", getTrackByIDWrapper(db))
}

type TrackExpanded struct {
	ID        uint     `json:"id"`
	Title     string   `json:"title"`
	Artists   []Artist `json:"artists"`
	Album     Album    `json:"album"`
	Duration  uint     `json:"duration"`
	SourceURL string   `json:"sourceUrl"`
}

type Track struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Duration  uint   `json:"duration"`
	SourceURL string `json:"sourceUrl"`
}

func DBTrackToAPITrack(track database.Track) Track {
	return Track{
		ID:        track.ID,
		Title:     track.Title,
		Duration:  track.Duration,
		SourceURL: fmt.Sprintf("%s/tracks/%d/audio.ogg", os.Getenv("ASSETS_SERVER_URI"), track.ID),
	}
}

func DBTrackExpandedToAPITrackExpanded(track database.TrackExpanded) TrackExpanded {
	artists := make([]Artist, len(track.Artists))
	for i, e := range track.Artists {
		artists[i] = DBArtistToAPIArtist(e)
	}

	return TrackExpanded{ID: track.ID,
		Title:     track.Title,
		Artists:   artists,
		Album:     DBAlbumToAPIAlbum(*track.Album),
		Duration:  track.Duration,
		SourceURL: fmt.Sprintf("%s/tracks/%d/audio.ogg", os.Getenv("ASSETS_SERVER_URI"), track.ID),
	}
}

func getTrackByIDWrapper(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id64, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Printf("[INFO] Non integer id \"%s\" has been provided: %s", idString, err)
			badRequestResponse(c, "Invalid track id.")
			return
		}
		id := uint(id64)
		track, err := database.GetTrackExpandedByID(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(c, fmt.Sprintf("Track with id \"%d\" was not found.", id))
				return
			}
			internalServerErrorResponse(c, "Internal server error.")
			return
		}
		c.JSON(http.StatusOK, DBTrackExpandedToAPITrackExpanded(*track))
	}
}
