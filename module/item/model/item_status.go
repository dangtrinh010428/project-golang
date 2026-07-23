package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
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
