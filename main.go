package main

import (
	"fmt"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TodoItem struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	// Image string `json:"image"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	// omitempty dùng để nếu value null thì không hiển thị
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func main() {
	fmt.Println("Hello, World!")
	now := time.Now().UTC()
	item := TodoItem{
		Id:          1,
		Title:       "Test example todo 22",
		Description: "Description todo",
		Status:      "Doing",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	r := gin.Default()

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": item,
			
		})
	})

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
