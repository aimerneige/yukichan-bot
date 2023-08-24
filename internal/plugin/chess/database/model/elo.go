package model

import "gorm.io/gorm"

// ELO user elo info
type ELO struct {
	gorm.Model
	Uin  int64 `gorm:"unique_index"`
	Name string
	Rate int
}
