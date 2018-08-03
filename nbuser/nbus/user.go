package nbus

import (
	"crypto/md5"
	"fmt"

	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
)

type UserService struct {
	BaseService
}

func NewUserService(ctx *context.Context) *UserService {
	return &UserService{BaseService{Ctx: ctx}}
}

// GetRequestUserType 获取用户登录方式
// 如果 http header 指定 nb-user-type，这优先按 header 指定的方式处理；服务端可配置默认的处理方式；兜底为 unkown
func (self *UserService) GetRequestUserType() string {
	userType := self.Ctx.Input.Header(nbum.HeaderUserType)
	defaultLoginUserType := beego.AppConfig.DefaultString("nbUserDefaultLoginType", nbum.UserTypeUnkown)

	switch userType {
	case "":
		return defaultLoginUserType
	case nbum.UserTypeEmail:
		return nbum.UserTypeEmail
	case nbum.UserTypeWX:
		return nbum.UserTypeWX
	default:
		return defaultLoginUserType
	}
}

func (self *UserService) GetUserByExternalInfo(exType, exId, exPassword string) (u *nbum.User, err error) {
	userId, err := nbum.GetUserIdByExternalId(exType, exId, self.encryptExPassword(exPassword))
	if err != nil {
		return nil, err
	}

	u, err = nbum.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	u.UserTypes, err = nbum.GetUserLoginTypesByUserId(userId)
	if err != nil {
		return nil, err
	}

	return
}

func (self *UserService) CreateNewUserWithLoginType(u *nbum.User, exType, exId, exPassword string) (err error) {
	o := orm.NewOrm()
	o.Begin()

	if err = u.CreateNew(o); err != nil {
		return err
	}

	if ult, err := self.BindNewLoginTypeForUser(u.Id, exType, exId, exPassword, o); err != nil {
		o.Rollback()
		return err
	} else {
		u.UserTypes = map[string]*nbum.UserLoginTypes{exType: ult}
	}

	o.Commit()
	return nil
}

func (self *UserService) BindNewLoginTypeForUser(userId int64, exType, exId, exPassword string, o ...orm.Ormer) (
	ult *nbum.UserLoginTypes, err error) {

	ult = &nbum.UserLoginTypes{
		UserId:     userId,
		ExType:     exType,
		ExId:       exId,
		ExPassword: self.encryptExPassword(exPassword),
	}

	return ult, ult.CreateNew(o...)
}

func (self *UserService) encryptExPassword(exPassword string) string {
	str := beego.AppConfig.String("nbUserEncryptToken") + ":" + exPassword
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
