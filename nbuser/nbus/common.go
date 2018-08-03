package nbus

import (
	"errors"
	"time"

	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbut"
	"github.com/NoneBorder/tasker"
	"github.com/astaxie/beego/context"
)

type CommonService struct {
	BaseService
}

func NewCommonService(ctx *context.Context) *CommonService {
	return &CommonService{BaseService{Ctx: ctx}}
}

func (self *CommonService) CreateNewVerifyCode(notifyType, notifyTo string) error {
	vc := self.getServerVerifyCode()
	if vc != nil && time.Now().Sub(vc.Created) < time.Minute {
		// 生成间隔至少1m
		return errors.New("verify code is generated too quickly")
	}
	if vc == nil {
		vc = new(nbum.VerifyCode)
	}

	vc.GenerateNew()
	if err := self.notifyVerifyCode(notifyType, notifyTo, vc); err != nil {
		return err
	}
	self.setServerVerifyCode(vc)
	return nil
}

func (self *CommonService) VerifyCode(code string) bool {
	vc := self.getServerVerifyCode()
	self.Logger().Debug().Str("clientCode", code).Interface("serverCode", vc).Msg("verify code debug")
	ret := vc != nil && vc.Verify(code)
	if ret {
		// 认证通过，删除验证码
		self.setServerVerifyCode(nil)
	}
	return ret
}

func (self *CommonService) notifyVerifyCode(notifyType, notifyTo string, vc *nbum.VerifyCode) error {
	sm := &nbut.SendMessage{
		Type:  notifyType,
		To:    notifyTo,
		Title: vc.NotifyTitle(),
		Body:  vc.NotifyBody(),
	}
	return tasker.MsgQPublishWithRetry(sm, 10*time.Second, 3)
}

func (self *CommonService) getServerVerifyCode() *nbum.VerifyCode {
	dI := self.Ctx.Input.Session(nbum.SessionKeyVerifyCode)
	if d, ok := dI.(*nbum.VerifyCode); ok {
		return d
	}
	return nil
}

func (self *CommonService) setServerVerifyCode(vc *nbum.VerifyCode) {
	self.Ctx.Output.Session(nbum.SessionKeyVerifyCode, vc)
}
