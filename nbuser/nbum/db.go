package nbum

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var ormDBAlias string = "default"
var ormDBTablePrefix string = ""

func DBInit(dbAlias, dbTablePrefix string) {
	if dbAlias != "" {
		ormDBAlias = dbAlias
	}
	if dbTablePrefix != "" {
		ormDBTablePrefix = dbTablePrefix
	}

	RegisterModel()
	go keepDBAlive()
}

func RegisterModel() {
	orm.RegisterModelWithPrefix(ormDBTablePrefix, new(User), new(UserLoginTypes), new(UserSignAuth))
}

func keepDBAlive() {
	o := orm.NewOrm()
	o.Using(ormDBAlias)
	for {
		o.Raw("select 1").Exec()
		time.Sleep(5 * time.Minute)
	}
}
