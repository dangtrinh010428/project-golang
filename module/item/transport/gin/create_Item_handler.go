package ginitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo.com/mod/common"
	"todo.com/mod/module/item/biz"
	"todo.com/mod/module/item/model"
	"todo.com/mod/module/item/storage"
)

// Tạo mới item
func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		// chỉ nhận body json
		var data model.TodoItemCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store := storage.NewSQLStore(db)
		business := biz.NewCreateItemBiz(store)

		if err := business.CreateNewItem(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
