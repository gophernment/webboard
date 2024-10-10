package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	// _ "modernc.org/sqlite" // windows
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
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, "SELECT name, message FROM webboard")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		tables := []board{}
		for rows.Next() {
			var name, message string
			if err := rows.Scan(&name, &message); err != nil {
				slog.Error(err.Error())
			}
			tables = append(tables, board{Name: name, Message: message})
		}
		c.JSON(http.StatusOK, tables)
	})
	r.POST("/api/boards", func(c *gin.Context) {
		var msg board
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		_, err = db.ExecContext(ctx, "INSERT INTO webboard(name,message) VALUES(?,?)", msg.Name, msg.Message)
		if err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		rows, err := db.QueryContext(ctx, "SELECT name, message FROM webboard")
		if err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		tables := []board{}
		for rows.Next() {
			var name, message string
			if err := rows.Scan(&name, &message); err != nil {
				slog.Error(err.Error())
			}
			tables = append(tables, board{Name: name, Message: message})
		}
		c.JSON(http.StatusOK, tables)
	})

	fmt.Println("listening and serving on :", os.Getenv("PORT"))
	r.Run()
	fmt.Println("bye")
}

type board struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
