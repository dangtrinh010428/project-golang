package biz

import (
	"context"
	"strings"

	"todo.com/mod/module/item/model"
)

type CreateItemStorage interface {
	createItem(ctx context.Context, data *model.TodoItemCreation) error
}

type createItemBiz struct {
	store CreateItemStorage
}

func NewCreateItemBiz(store CreateItemStorage) *createItemBiz {
	return &CreateItemStorage{store}
}

func (biz *createItemBiz) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	title := strings.TrimSpace(data.Title)
	if title == "" {
		return model.ErrTiteIsBlank
	}

	if err := biz.store.createItem(ctx, data), err != nil {
		return err
	}
	return nil
}
