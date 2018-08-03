package nbuc

import (
	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbus"
	"github.com/NoneBorder/bee-modules/nbuser/nbutils"
	"github.com/NoneBorder/dora/apiresp"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/rs/zerolog"
)

type BaseController struct {
	beego.Controller
	logger *zerolog.Logger
}

func (self *BaseController) ResetLogger() {
	self.logger = nbutils.BaseLoggerSetup(self.Ctx, func(c zerolog.Context) zerolog.Context {
		if u := nbus.NewLoginStateService(self.Ctx).JustGetUser(); u != nil {
			c = c.Int64("userId", u.Id)
		} else {
			c = c.Int64("userId", 0)
		}
		return c
	})
}

func (self *BaseController) Logger(refresh ...bool) *zerolog.Logger {
	if self.logger == nil {
		self.ResetLogger()
	}
	return self.logger
}

func (self *BaseController) ValidResWrapper(ret bool, err error, valid *validation.Validation) {
	if err != nil {
		self.Logger().Error().Err(err).Msg("valid got error")
		apiresp.NewDetail(nbum.ErrCodeUnkownError, nbum.ErrUnkownError.Error()).JSON(self.Ctx)
	} else if !ret {
		// valid 没有通过
		errs := make([]map[string]string, len(valid.Errors))
		for i, e := range valid.Errors {
			errs[i] = map[string]string{
				"Field":   e.Field,
				"Message": e.Message,
			}
		}
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "invalid req data", errs).JSON(self.Ctx)
	}
}
