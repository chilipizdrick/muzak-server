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

const ARTISTS_API_ROUTE = "/artists"

func assignArtistsRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(ARTISTS_API_ROUTE)

	group.GET("/:id", getArtistByIDWrapper(db))
}

type Artist struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsVerified bool   `json:"isVerified"`
}

type ArtistExpanded struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Albums     []Album `json:"albums"`
	Tracks     []Track `json:"tracks"`
	IsVerified bool    `json:"isVerified"`
}

func DBArtistToAPIArtist(artist database.Artist) Artist {
	return Artist{
		ID:         artist.ID,
		Name:       artist.Name,
		IsVerified: artist.IsVerified,
	}
}

func DBArtistExpandedToAPIArtistExpanded(artist database.ArtistExpanded) ArtistExpanded {
	albums := make([]Album, len(artist.Albums))
	for i, e := range artist.Albums {
		albums[i] = DBAlbumToAPIAlbum(e)
	}

	tracks := make([]Track, len(artist.Tracks))
	for i, e := range artist.Tracks {
		tracks[i] = DBTrackToAPITrack(e)
	}

	return ArtistExpanded{
		ID:         artist.ID,
		Name:       artist.Name,
		Albums:     albums,
		Tracks:     tracks,
		IsVerified: artist.IsVerified,
	}
}

func getArtistByIDWrapper(db *gorm.DB) func(*gin.Context) {
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
		artist, err := database.GetArtistExpandedByID(db, id)
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
		c.IndentedJSON(http.StatusOK, DBArtistExpandedToAPIArtistExpanded(*artist))
	}
}
