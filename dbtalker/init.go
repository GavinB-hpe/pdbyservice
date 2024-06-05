package dbtalker

import "gorm.io/gorm"

type DBTalker struct {
	DB *gorm.DB
}

func NewDBTalker(db *gorm.DB) *DBTalker {
	return &DBTalker{
		DB: db,
	}
}
