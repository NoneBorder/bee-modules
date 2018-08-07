package nbuc

import (
	"encoding/json"

	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbus"
	"github.com/NoneBorder/dora/apiresp"
)

type UserController struct {
	BaseController
}

func (self *UserController) UserName() {
	u := nbus.NewLoginStateService(self.Ctx).GetUser()
	self.Ctx.WriteString(u.Name)
}

func (self *UserController) WXLoginCallback() {
	code := self.GetString("code")
	if code == "" {
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "login failed").JSON(self.Ctx)
	}

	ls := nbus.NewLoginStateService(self.Ctx)

	if err := ls.WXInitAccessTokenByCode(code); err != nil {
		self.Logger().Error().Err(err).Msg("init access token failed")
		errRedirect := self.GetString(nbum.AuthLoginWXCallbackErrKey, "")
		if errRedirect == "" {
			apiresp.NewDetail(nbum.ErrCodeGeneralError, "init access token failed").JSON(self.Ctx)
		}
		self.Redirect(errRedirect, 302)
		return
	}

	succRedirect := self.GetString(nbum.AuthLoginWXCallbackSuccKey, "")
	if succRedirect == "" {
		self.Logger().Error().Msg("not set success redirect url")
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "not set success redirect url").JSON(self.Ctx)
	}

	self.Redirect(succRedirect, 302)
}

type RegisterForm struct {
	ExType     string
	ExId       string
	ExPassword string
	User       *nbum.User
	VerifyCode string
}

func (self *UserController) Register() {
	self.Logger().Debug().Bytes("requestBody", self.Ctx.Input.RequestBody).Msg("request body")

	registerForm := RegisterForm{}
	if err := json.Unmarshal(self.Ctx.Input.RequestBody, &registerForm); err != nil {
		self.Logger().Error().Err(err).Msg("parse request body failed")
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "parse request body failed").JSON(self.Ctx)
	}
	if registerForm.User == nil {
		registerForm.User = new(nbum.User)
	}
	if registerForm.User.Name == "" {
		registerForm.User.Name = registerForm.ExId
	}
	self.Logger().Debug().Interface("registerForm", registerForm).Msg("register form")

	// 当前仅支持邮箱注册
	if registerForm.ExType != nbum.UserTypeEmail {
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "invalid ExType").JSON(self.Ctx)
	}
	if registerForm.ExId == "" || registerForm.ExPassword == "" {
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "empty Id or Password").JSON(self.Ctx)
	}

	// 验证码认证
	commonService := nbus.NewCommonService(self.Ctx)
	if !commonService.VerifyCode(registerForm.VerifyCode) {
		// 没有通过验证码认证
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "invalid verify code").JSON(self.Ctx)
	}

	userService := nbus.NewUserService(self.Ctx)
	err := userService.CreateNewUserWithLoginType(
		registerForm.User, registerForm.ExType, registerForm.ExId, registerForm.ExPassword)

	if err != nil {
		self.Logger().Error().Err(err).Msg("register new user failed")
		apiresp.NewDetail(nbum.ErrCodeGeneralError, "register failed:"+err.Error()).JSON(self.Ctx)
	}

	apiresp.NewResp(registerForm.User).JSON(self.Ctx)
}

func (self *UserController) Login() {
	self.Logger().Debug().Bytes("requestBody", self.Ctx.Input.RequestBody).Msg("request body")

	loginForm := RegisterForm{}
	if err := json.Unmarshal(self.Ctx.Input.RequestBody, &loginForm); err != nil {
		self.Logger().Error().Err(err).Msg("parse request body failed")
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "parse request body failed").JSON(self.Ctx)
	}
	self.Logger().Debug().Interface("loginForm", loginForm).Msg("login form")

	if loginForm.ExId == "" || loginForm.ExPassword == "" {
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "empty Id or Password").JSON(self.Ctx)
	}

	switch loginForm.ExType {
	case nbum.UserTypeEmail:
		ls := nbus.NewLoginStateService(self.Ctx)
		user, err := ls.EmailLogin(loginForm.ExId, loginForm.ExPassword)
		if err != nil {
			self.Logger().Error().Err(err).Msg("login failed")
			apiresp.NewDetail(nbum.ErrCodeBadRequest, "login failed:"+err.Error()).JSON(self.Ctx)
		}
		apiresp.NewResp(user).JSON(self.Ctx)

	default:
		// 当前仅支持邮箱注册
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "invalid ExType").JSON(self.Ctx)
		return
	}
}

func (self *UserController) Logout() {
	nbus.NewLoginStateService(self.Ctx).Logout()
	apiresp.NewResp(nil).JSON(self.Ctx)
}
