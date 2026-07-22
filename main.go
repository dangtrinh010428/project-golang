package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

type TodoItemCreation struct {
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status"`
}

func (TodoItemCreation) TableName() string {
	return "todo_items"
}

func main() {

	dsn := os.Getenv("MYSQL_CONNECTION")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Lỗi connect database: ", err)
	}

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

	// CRUD: Create Read Update Delete
	// POST /v1/items create new
	// GET /v1/items list items
	// GET /v1/item/:id get item by ID
	// PUT thay đổi cả object
	// PATCH thay đổi từng thành phần
	// POST  /v1/items/:id update item by ID
	// DELETE /v1/items/:id Delete item by ID

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItem(db))
			items.GET("")
			items.GET("/:id")
			items.PATCH("/:id")
			items.DELETE("/:id")
		}
	}

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

func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		// chỉ nhận body json
		var data TodoItemCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
