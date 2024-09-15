package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

type OKResponse struct {
	OK OK `json:"ok"`
}

type OK struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

func badRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: Error{
			Status:  http.StatusBadRequest,
			Message: message,
		},
	})
}

func internalServerErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: Error{
			Status:  http.StatusInternalServerError,
			Message: message,
		},
	})
}

func notFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: Error{
			Status:  http.StatusNotFound,
			Message: message,
		},
	})
}
