package structs

import (
	"database/sql/driver"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	return "{" + strings.Join(s, ",") + "}", nil
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	val, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}
	*s = strings.Split(strings.Trim(val, "{}"), ",")
	return nil
}

type Food struct {
	gorm.Model
	CategoryID  uint        `gorm:"not null"`
	Name        string      `gorm:"size:255;not null"`
	Description string      `gorm:"size:255"`
	Quantity    uint        `gorm:"not null"`
	Pictures    StringArray `gorm:"type:text[]"`
	Category    Category    `gorm:"foreignKey:CategoryID"`
}
