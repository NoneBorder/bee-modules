package nbus

import (
	"encoding/gob"

	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbutils"
	"github.com/astaxie/beego/context"
	"github.com/rs/zerolog"
	"github.com/silenceper/wechat/oauth"
)

func init() {
	gob.Register(new(nbum.User))
	gob.Register(new(nbum.VerifyCode))
	gob.Register(oauth.ResAccessToken{})
}

type BaseService struct {
	Ctx    *context.Context
	logger *zerolog.Logger
}

func (self *BaseService) ResetLogger() {
	self.logger = nbutils.BaseLoggerSetup(self.Ctx, func(c zerolog.Context) zerolog.Context {
		if u := NewLoginStateService(self.Ctx).JustGetUser(); u != nil {
			c = c.Int64("userId", u.Id)
		} else {
			c = c.Int64("userId", 0)
		}
		return c
	})
}

func (self *BaseService) Logger() *zerolog.Logger {
	if self.logger == nil {
		self.ResetLogger()
	}
	return self.logger
}
