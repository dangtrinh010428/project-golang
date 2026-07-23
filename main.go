package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"todo.com/mod/common"
)

// ENUM
type ItemStatus int

const (
	ItemStatusDoing ItemStatus = iota
	ItemStatusDone
	ItemStatusDeleted
)

var AllItemStatus = [3]string{"Doing", "Done", "Deleted"}

func (item *ItemStatus) String() string {
	return AllItemStatus[*item]
}

func parseStr2ItemStatus(s string) (ItemStatus, error) {
	for i := range AllItemStatus {
		if AllItemStatus[i] == s {
			return ItemStatus(i), nil
		}

	}
	return ItemStatus(0), errors.New("Invalid item status")
}

func (item *ItemStatus) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return errors.New(fmt.Sprintln("Fail to scan data from sql", value))
	}

	v, err := parseStr2ItemStatus(string(bytes))

	if err != nil {
		return errors.New(fmt.Sprintln("Fail to scan data from sql", value))
	}

	*item = v
	return nil
}

func (item *ItemStatus) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return item.String(), nil
}

func (item *ItemStatus) MarshalJSON() ([]byte, error) {
	if item == nil {
		return nil, nil
	}
	return []byte(fmt.Sprintf("\"%s\"", item.String())), nil
}

func (item *ItemStatus) UnMarshalJSON(data []byte) error {
	str := strings.ReplaceAll(string(data), "\"", "")
	itemValue, err := parseStr2ItemStatus(str)
	if err != nil {
		return err
	}

	*item = itemValue
	return nil
}

type TodoItem struct {
	// kỹ thuật embeding struct
	common.SqlModel
	Title string `json:"title" gorm:"column:title"`
	// Image string `json:"image"`
	Description string      `json:"description" gorm:"column:description"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Id          int         `json:"-" gorm:"column:id"`
	Title       string      `json:"title" gorm:"column:title;"`
	Description string      `json:"description" gorm:"column:description;"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdate struct {
	Title       *string     `json:"title" gorm:"column:title;"`
	Description *string     `json:"description" gorm:"column:description;"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`
}

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()
}

func main() {

	dsn := os.Getenv("MYSQL_CONNECTION")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Lỗi connect database: ", err)
	}
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreateItem(db))
			items.GET("", GetAllItem(db))
			items.GET("/:id", GetItem(db))
			items.PATCH("/:id", Update(db))
			items.DELETE("/:id", Delete(db))
		}
	}

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// Tạo mới item
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

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		// chỉ nhận body json
		var data TodoItem

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}

func GetAllItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		paging.Process()

		var result []TodoItem

		if err := db.Table(TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		db := db.Where("status <> ?", "Deleted")

		if err := db.Order("id desc").Offset((paging.Page - 1) * paging.Limit).Limit(paging.Limit).Find(&result).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, nil))
	}
}

func Update(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemUpdate
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func Delete(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{"status": "Deleted"}).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
