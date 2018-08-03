package nbum

import (
	"crypto/md5"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserSignAuth struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"userId" orm:"unique"`
	SecretKey string    `json:"secretKey"`
	Created   time.Time `json:"created" orm:"auto_now_add;type(datetime)"`
	Updated   time.Time `json:"updated" orm:"auto_now;type(datetime)"`
}

func GetUserSignByUserId(userId int64) (usa *UserSignAuth, err error) {
	o := orm.NewOrm()
	o.Using(ormDBAlias)
	usa = new(UserSignAuth)
	err = o.QueryTable(new(UserSignAuth)).Filter("UserId", userId).One(usa)
	return
}

func (self *UserSignAuth) Auth(params map[string]string, rawSign string) (*User, error) {
	if rawSign == "" || self.UserId == 0 {
		return nil, errors.New("invalid sign auth input")
	}

	sign := self.Sign(params)
	if sign != rawSign {
		return nil, errors.New("sign auth failed")
	}

	u, err := GetUserById(self.UserId)
	return u, err
}

func (self *UserSignAuth) Sign(params map[string]string) string {
	if self.SecretKey == "" {
		return ""
	}

	params[SignParamSecretKey] = self.SecretKey

	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	items := make([]string, len(keys))
	for i, k := range keys {
		items[i] = k + "=" + params[k]
	}
	delete(params, SignParamSecretKey)

	itemsStr := strings.Join(items, "||")
	sign := fmt.Sprintf("%x", md5.Sum([]byte(itemsStr)))
	//zlog.Debug().Str("rawStr", itemsStr).Str("sign", sign).Msg("sign info")
	return sign
}
