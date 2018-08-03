package nbum

import "github.com/astaxie/beego/orm"

type UserLoginTypes struct {
	Id         int64
	UserId     int64  `json:"userId"`
	ExType     string `json:"extype"`
	ExId       string `json:"exid"`
	ExPassword string `json:"null"`
}

// 多字段唯一键
func (self *UserLoginTypes) TableUnique() [][]string {
	return [][]string{
		[]string{"ExType", "ExId"},
	}
}

// 设置引擎为 INNODB
func (self *UserLoginTypes) TableEngine() string {
	return "INNODB"
}

func GetUserIdByExternalId(exType, exID, exPassword string) (userId int64, err error) {
	o := orm.NewOrm()
	o.Using(ormDBAlias)
	ult := new(UserLoginTypes)
	qs := o.QueryTable(ult).Filter("ExType", exType).Filter("ExId", exID)
	if exType == UserTypeEmail {
		// 如果是 邮箱用户，需要验证密码
		qs = qs.Filter("ExPassword", exPassword)
	}
	err = qs.One(ult)
	return ult.UserId, err
}

func GetUserLoginTypesByUserId(userId int64) (ret map[string]*UserLoginTypes, err error) {
	o := orm.NewOrm()
	o.Using(ormDBAlias)

	var ults []*UserLoginTypes
	if _, err = o.QueryTable(new(UserLoginTypes)).Filter("UserId", userId).All(&ults); err != nil {
		return nil, err
	}

	ret = map[string]*UserLoginTypes{}
	for _, ult := range ults {
		ret[ult.ExType] = ult
	}

	return
}

func (self *UserLoginTypes) CreateNew(o ...orm.Ormer) error {
	o = append(o, orm.NewOrm())
	o[0].Using(ormDBAlias)

	_, err := o[0].Insert(self)
	return err
}
