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

const USERS_API_ROUTE = "/users"

func assignUsersRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(USERS_API_ROUTE)

	group.GET("/:id", getUserByIDWrapper(db))
}

type User struct {
	ID                  uint         `json:"id"`
	Username            string       `json:"username"`
	PlayerState         *PlayerState `json:"playerState"`
	IsPlayerStatePublic bool         `json:"isPlayerStatePublic"`
}

type UserExpanded struct {
	ID                  uint         `json:"id"`
	Username            string       `json:"username"`
	Playlists           []Playlist   `json:"playlists"`
	PlayerState         *PlayerState `json:"playerState"`
	IsPlayerStatePublic bool         `json:"isPlayerStatePublic"`
}

type PlayerState struct {
	Track                   Track   `json:"track"`
	IsPlaying               bool    `json:"isPlaying"`
	Progress                uint    `json:"progress"` // In seconds
	Device                  string  `json:"device"`
	IsShuffleEnabled        bool    `json:"isShuffleEnabled"`
	IsRepeatPlaylistEnabled bool    `json:"isRepeatPlaylistEnabled"`
	IsRepeatTrackEnabled    bool    `json:"isRepeatTrackEnabled"`
	Volume                  float64 `json:"volume"` // From 0.0 to 1.0
}

func DBUserToAPIUser(user database.User) User {
	playerState := new(PlayerState)
	if user.IsPlayerStatePublic {
		*playerState = DBPlayerStateToAPIPlayerState(user.PlayerState)
	} else {
		playerState = nil
	}

	return User{
		ID:                  user.ID,
		Username:            user.Username,
		PlayerState:         playerState,
		IsPlayerStatePublic: user.IsPlayerStatePublic,
	}
}

func DBPlayerStateToAPIPlayerState(playerState *database.PlayerState) PlayerState {
	return PlayerState{
		Track:                   DBTrackToAPITrack(playerState.Track),
		IsPlaying:               playerState.IsPlaying,
		Progress:                playerState.Progress,
		Device:                  playerState.Device,
		IsShuffleEnabled:        playerState.IsShuffleEnabled,
		IsRepeatPlaylistEnabled: playerState.IsRepeatPlaylistEnabled,
		IsRepeatTrackEnabled:    playerState.IsRepeatTrackEnabled,
		Volume:                  playerState.Volume,
	}
}

func DBUserExpandedToAPIUserExpanded(user database.UserExpanded) UserExpanded {
	playlists := make([]Playlist, len(user.Playlists))
	for i, e := range user.Playlists {
		playlists[i] = DBPlaylistToAPIPlaylist(e)
	}

	playerState := new(PlayerState)
	if user.IsPlayerStatePublic {
		*playerState = DBPlayerStateToAPIPlayerState(user.PlayerState)
	} else {
		playerState = nil
	}

	return UserExpanded{
		ID:                  user.ID,
		Username:            user.Username,
		Playlists:           playlists,
		PlayerState:         playerState,
		IsPlayerStatePublic: user.IsPlayerStatePublic,
	}
}

func getUserByIDWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id64, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Printf("[INFO] Non integer id \"%s\" has been provided: %s", idString, err)
			c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{
				Error: Error{
					Status:  http.StatusInternalServerError,
					Message: "Invalid user id.",
				},
			})
			return
		}
		id := uint(id64)
		user, err := database.GetUserExpandedByID(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.IndentedJSON(http.StatusNotFound, ErrorResponse{
					Error: Error{
						Status:  http.StatusNotFound,
						Message: fmt.Sprintf("User with id \"%d\" was not found.", id),
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
		c.IndentedJSON(http.StatusOK, DBUserExpandedToAPIUserExpanded(*user))
	}
}
