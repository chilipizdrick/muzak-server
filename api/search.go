package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const SEARCH_API_ROUTE = "/search"

func assignSearchRouteHandlers(parentGroup *gin.RouterGroup, db *gorm.DB) {
	group := parentGroup.Group(SEARCH_API_ROUTE)

    group.GET("/:type", searchWrapper(db))
}

type SearchRequest struct {
	Type  string `validate:"required,oneof=album artist playlist track user"`
	Query string `validate:"required,len>0"`
}

func searchWrapper(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		request := SearchRequest{
			Type:  c.Param("type"),
			Query: c.Param("q"),
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		if err := validate.Struct(request); err != nil {
			log.Printf("[INFO] Failed to validate search request: %s", err)
			c.IndentedJSON(http.StatusInternalServerError, ErrorResponse{
				Error: Error{
					Status:  http.StatusInternalServerError,
					Message: "Invalid search request.",
				},
			})
			return
		}

		switch request.Type {
		case "album":
		case "artist":
		case "playlist":
		case "track":
		case "user":
		default:
			panic("unreachable")
		}
	}
}
