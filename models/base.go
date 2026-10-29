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
// 注：mattn 依赖 CGO，编译需 CGO_ENABLED=1。纯 Go 的 modernc.org/sqlite 虽可 CGO_ENABLED=0
// 交叉编译，但纯 Go 实现比 C 库慢数倍（实测任务执行 6s → 45s），故回退 mattn。
func InitSQLite(path string) error {
	var e error
	Orm, e = gorm.Open("sqlite3", path)
	return e
}
