package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./webboard.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.StaticFile("/webboard.html", "./webboard.html")

	r.GET("/api/boards", func(c *gin.Context) {
		c.JSON(http.StatusOK, webboard)
	})
	r.POST("/api/boards", func(c *gin.Context) {
		var msg board
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		webboard = append(webboard, msg)

		c.JSON(http.StatusOK, webboard)
	})

	fmt.Println("listening and serving on :", os.Getenv("PORT"))
	r.Run()
	fmt.Println("bye")
}

type board struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var webboard = []board{}
