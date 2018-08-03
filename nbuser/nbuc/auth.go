package nbuc

import (
	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbus"
	"github.com/NoneBorder/dora/apiresp"
)

type AuthController struct {
	BaseController
}

func (self *AuthController) UserName() {
	u := nbus.NewLoginStateService(self.Ctx).GetUser()
	self.Ctx.WriteString(u.Name)
}

func (self *AuthController) WXLoginCallback() {
	code := self.GetString("code")
	if code == "" {
		apiresp.NewDetail(400, "login failed").BeegoServeJSON(self.Controller)
	}

	ls := nbus.NewLoginStateService(self.Ctx)

	if err := ls.WXInitAccessTokenByCode(code); err != nil {
		self.Logger().Error().Err(err).Msg("init access token failed")
		errRedirect := self.GetString(nbum.AuthLoginWXCallbackErrKey, "")
		if errRedirect == "" {
			apiresp.NewDetail(500, "init access token failed")
			return
		}
		self.Redirect(errRedirect, 302)
		return
	}

	succRedirect := self.GetString(nbum.AuthLoginWXCallbackSuccKey, "")
	if succRedirect == "" {
		self.Logger().Error().Msg("not set success redirect url")
		apiresp.NewDetail(400, "not set success redirect url")
		return
	}

	self.Redirect(succRedirect, 302)
}
