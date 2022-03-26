package infrastructure

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func NewMySQL(conString string) *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(conString), &gorm.Config{
			PrepareStmt: true,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
	)

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database : %s", err.Error()))
	}

	return db
}
