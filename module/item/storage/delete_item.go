package storage

import (
	"context"

	"todo.com/mod/module/item/model"
)

func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {

	deletedStatus := model.ItemStatusDeleted

	if err := s.db.Table(model.TodoItem{}.TableName()).
		Where(cond).
		Updates(map[string]interface{}{"status": deletedStatus.String()}).Error; err != nil {
		return err
	}
	return nil
}
