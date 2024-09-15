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

const ALBUMS_API_ROUTE = "/albums"

func assignAlbumsRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(ALBUMS_API_ROUTE)

	group.GET("/:id", getAlbumByIDWrapper(db))
}

type Album struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type AlbumExpanded struct {
	ID      uint     `json:"id"`
	Title   string   `json:"title"`
	Artists []Artist `json:"artists"`
	Tracks  []Track  `json:"tracks"`
}

func DBAlbumToAPIAlbum(album database.Album) Album {
	return Album{
		ID:    album.ID,
		Title: album.Title,
	}
}

func DBAlbumExpandedToAPIAlbumExpanded(album database.AlbumExpanded) AlbumExpanded {
	artists := make([]Artist, len(album.Artists))
	for i, e := range album.Artists {
		artists[i] = DBArtistToAPIArtist(e)
	}

	tracks := make([]Track, len(album.Tracks))
	for i, e := range album.Tracks {
		tracks[i] = DBTrackToAPITrack(e)
	}

	return AlbumExpanded{
		ID:      album.ID,
		Title:   album.Title,
		Artists: artists,
		Tracks:  tracks,
	}
}

func getAlbumByIDWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id64, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Printf("[INFO] Non integer id \"%s\" has been provided: %s", idString, err)
			badRequestResponse(c, "Invalid album id.")
			return
		}
		id := uint(id64)
		album, err := database.GetAlbumExpandedByID(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				notFoundResponse(c, fmt.Sprintf("Album with id \"%d\" was not found.", id))
				return
			}
			internalServerErrorResponse(c, "Internal server error.")
			return
		}
		c.JSON(http.StatusOK, DBAlbumExpandedToAPIAlbumExpanded(*album))
	}
}
