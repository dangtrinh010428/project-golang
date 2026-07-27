package model

import (
	"errors"

	"todo.com/mod/common"
)

var (
	ErrTiteIsBlank = errors.New("Title cannot be blank")
	ErrItemDeleted = errors.New("Item is deleted")
)

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

// ENUM
type ItemStatus int

const (
	ItemStatusDoing ItemStatus = iota
	ItemStatusDone
	ItemStatusDeleted
)
