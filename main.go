package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Track struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	// Cover  image.Image `json:"cover"`
}

var Tracks = []Track{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Album: "a"},
	{ID: 2, Title: "Blue Train", Artist: "John Coltrane", Album: "a"},
	{ID: 3, Title: "Blue Train", Artist: "John Coltrane", Album: "a"},
	{ID: 4, Title: "Blue Train", Artist: "John Coltrane", Album: "a"},
}

func getTracks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Tracks)
}

func main() {
	router := gin.Default()
	router.GET("/tracks", getTracks)
	router.Run("localhost:8080")
}