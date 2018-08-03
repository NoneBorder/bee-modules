package nbum

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id         int64                      `json:"id"`
	Name       string                     `json:"name"`
	Sex        int32                      `json:"sex"` //1 男性； 2 女性； 0 未知
	Province   string                     `json:"province"`
	City       string                     `json:"city"`
	Country    string                     `json:"country"`
	HeadImgURL string                     `json:"headimgurl" orm:"column(headimgurl)"`
	Created    time.Time                  `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated    time.Time                  `json:"updated" orm:"auto_now;type(datetime)"`
	UserTypes  map[string]*UserLoginTypes `json:"userTypes" orm:"-"`
}

// 设置引擎为 INNODB
func (self *User) TableEngine() string {
	return "INNODB"
}

func (self *User) CreateNew(o ...orm.Ormer) error {
	o = append(o, orm.NewOrm())
	o[0].Using(ormDBAlias)
	_, err := o[0].Insert(self)
	return err
}

func GetUserById(userId int64) (u *User, err error) {
	o := orm.NewOrm()
	o.Using(ormDBAlias)
	u = new(User)
	err = o.QueryTable(new(User)).Filter("Id", userId).One(u)
	return
}
