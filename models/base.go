package models

import (
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	Orm *gorm.DB
	err error
)

func Init(c db.Connection) {
	Orm, err = gorm.Open("mysql", c.GetDB("default"))

	if err != nil {
		panic("initialize orm failed")
	}
}

// InitSQLite 以嵌入式 SQLite 打开本地库（CLI 离线模式）。
// 使用 gorm 的 sqlite3 方言，底层驱动为 mattn/go-sqlite3（CGO）。
func InitSQLite(path string) error {
	var e error
	Orm, e = gorm.Open("sqlite3", path)
	return e
}
