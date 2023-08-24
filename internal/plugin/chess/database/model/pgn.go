package model

import "gorm.io/gorm"

// PGN chess pgn info
type PGN struct {
	gorm.Model
	Data      string
	WhiteUin  int64
	BlackUin  int64
	WhiteName string
	BlackName string
}
