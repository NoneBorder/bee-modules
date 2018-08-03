package nbus

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/NoneBorder/bee-modules/doralog"
	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/dora/apiresp"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"

	"github.com/NoneBorder/dora"
	"github.com/astaxie/beego"
	"github.com/silenceper/wechat/oauth"
)

/*

LoginStateService 登录状态管理

主要登录逻辑：
1. GetUser： 登录状态检测，如果当前未登录，直接截断请求，并返回 http code 401；当请求头包含明确的登录方式 header 'nb-user-type'
   字段时，按具体的登录方式返回 http body 内容；否则 body 内容为空

*/
type LoginStateService struct {
	BaseService
	userService *UserService
}

func NewLoginStateService(ctx *context.Context) *LoginStateService {
	return &LoginStateService{BaseService: BaseService{Ctx: ctx}, userService: NewUserService(ctx)}
}

func (self *LoginStateService) JustGetUser() *nbum.User {
	if u := self.signAuth(); u != nil {
		return u
	}

	uI := self.Ctx.Input.Session(nbum.SessionKeyUser)
	if u, ok := uI.(*nbum.User); ok {
		return u
	}
	return nil
}

func (self *LoginStateService) GetUser() *nbum.User {
	if u := self.signAuth(); u != nil {
		return u
	}

	return self.checkLogin()
}

func (self *LoginStateService) IsSignReq() bool {
	return self.Ctx.Input.Header(nbum.HeaderNBSign) != ""
}

func (self *LoginStateService) signAuth() *nbum.User {
	uI := self.Ctx.Input.GetData(nbum.CtxInputDataKeyUser)
	if u, ok := uI.(*nbum.User); ok {
		return u
	}

	var u *nbum.User = nil
	header := self.Ctx.Input.Header(nbum.HeaderNBSign)
	switch header {
	case nbum.NBSignTypeMd5Sign:
		u = self.md5SignAuth()
	}

	if u != nil {
		self.Ctx.Input.SetData(nbum.CtxInputDataKeyUser, u)
	}
	return u
}

func (self *LoginStateService) md5SignAuth() *nbum.User {
	// ts 认证
	ts := self.Ctx.Input.Query(nbum.SignParamTS)
	if ts == "" {
		return nil
	}
	tsInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil
	}
	if tsInt < time.Now().Unix()-300 {
		return nil
	}

	// userId 合法性检查
	userIdStr := self.Ctx.Input.Query(nbum.SignParamUserId)
	if userIdStr == "" {
		return nil
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil || userId == 0 {
		return nil
	}

	// sign 检查
	sign := self.Ctx.Input.Query(nbum.SignParamSign)
	if sign == "" {
		return nil
	}

	// 计算 sign
	pm := make(map[string]string)
	for k := range self.Ctx.Request.Form {
		pm[k] = self.Ctx.Request.Form.Get(k)
	}

	// request body
	if len(self.Ctx.Input.RequestBody) > 0 {
		pm[nbum.SignParamRequestBody] = string(self.Ctx.Input.RequestBody)
	}

	doralog.H(self.Ctx).Debug().Interface("pm", pm).Msg("sign auth raw data")

	delete(pm, nbum.SignParamSign)
	usa, err := nbum.GetUserSignByUserId(userId)
	if err != nil {
		doralog.H(self.Ctx).Error().Err(err).Msg("sign auth get secretkey failed")
		return nil
	}

	u, err := usa.Auth(pm, sign)
	if err != nil {
		doralog.H(self.Ctx).Error().Err(err).Msg("sign auth failed")
		return nil
	}

	return u
}

func (self *LoginStateService) WXInitAccessTokenByCode(code string) error {
	oauth := nbum.DefaultWX.Wechat.GetOauth()
	resToken, err := oauth.GetUserAccessToken(code)
	if err != nil {
		self.Logger().Error().Str("code", code).Err(err).Msg("get user access token by code failed")
		return fmt.Errorf("get user access token by code %s failed", code)
	}

	self.Logger().Debug().Interface("accessToken", resToken).Msg("get access token from wechat")
	self.Ctx.Output.Session(nbum.SessionKeyWXResAccessToken, resToken)

	user, err := self.getUserByWXToken(resToken)
	if err != nil {
		self.Logger().Error().Str("openid", resToken.OpenID).Err(err).Msg("get user by wechat token failed")
		return errors.New("get user by wechat token failed")
	}

	self.Logger().Debug().Interface("user", user).Str("openid", resToken.OpenID).Msg("get user by wechat token")
	self.Ctx.Output.Session(nbum.SessionKeyUser, user)
	return nil
}

func (self *LoginStateService) getUserByWXToken(resToken oauth.ResAccessToken) (*nbum.User, error) {
	user, err := self.userService.GetUserByExternalInfo(nbum.UserTypeWX, resToken.OpenID, "")
	if err == orm.ErrNoRows {
		// 不存在用户，开始新建
		return self.createNewUserByWXAccessToken(resToken)
	}

	if err != nil {
		self.Logger().Error().Err(err).Msg("get user from model failed")
		return nil, errors.New("get user from model failed")
	}

	return user, nil
}

func (self *LoginStateService) createNewUserByWXAccessToken(resToken oauth.ResAccessToken) (*nbum.User, error) {
	self.Logger().Info().Str("openid", resToken.OpenID).Msg("prepare to create new user")
	oauth := nbum.DefaultWX.Wechat.GetOauth()
	userInfo, err := oauth.GetUserInfo(resToken.AccessToken, resToken.OpenID)
	if err != nil {
		self.Logger().Error().Err(err).Msg("get userinfo from wechat failed")
		return nil, errors.New("get userinfo from wechat failed")
	}

	u := &nbum.User{
		Name:       userInfo.Nickname,
		Sex:        userInfo.Sex,
		Province:   userInfo.Province,
		City:       userInfo.City,
		Country:    userInfo.Country,
		HeadImgURL: userInfo.HeadImgURL,
	}

	err = self.userService.CreateNewUserWithLoginType(u, nbum.UserTypeWX, resToken.OpenID, "")
	if err != nil {
		self.Logger().Error().Err(err).Msg("create new user failed")
		return nil, errors.New("create new user failed")
	}

	self.Logger().Debug().Interface("user", u).Str("openid", resToken.OpenID).Msg("create new user success")
	return u, nil
}

func (self *LoginStateService) checkLogin() *nbum.User {
	uI := self.Ctx.Input.Session(nbum.SessionKeyUser)

	u, ok := uI.(*nbum.User)
	if !ok || u == nil {
		switch self.userService.GetRequestUserType() {
		case nbum.UserTypeWX:
			u = self.wxLogin()
		case nbum.UserTypeUnkown:
			u = self.unkownLogin()
		case "":
			self.Logger().Fatal().Msg("unknown user type")
		}

		// 成功重新登录后，重置 server logger
		if u != nil {
		}
	}
	return u
}

func (self *LoginStateService) Logout() {
	self.clearSession()
}

func (self *LoginStateService) EmailLogin(exId, exPassword string) (*nbum.User, error) {
	if user, err := self.userService.GetUserByExternalInfo(nbum.UserTypeEmail, exId, exPassword); err != nil {
		return nil, err
	} else {
		self.Ctx.Output.Session(nbum.SessionKeyUser, user)
		return user, nil
	}
}

// unkownLogin 不确定当前具体的登录方式，但用户又没有登录
// 处理方式为分直接请求和 ajax 两种：
// 1. ajax：直接截断请求，返回 http 401 状态码, http body 为空
// 2. 直接请求：重定向到用户登录页，如未配置按 ajax 方式返回
func (self *LoginStateService) unkownLogin() *nbum.User {
	if self.IsSignReq() || self.Ctx.Input.IsAjax() || beego.AppConfig.String("nbUserUnloginRedirect") == "" {
		// ajax
		resp := apiresp.NewDetail(nbum.ErrCodeNotLogin, "unauthorized request")
		resp.JSON(self.Ctx, 401)
	} else {
		self.Ctx.Redirect(302, beego.AppConfig.String("nbUserUnloginRedirect"))
	}
	panic(beego.ErrAbort)
	return nil
}

func (self *LoginStateService) wxLogin() *nbum.User {
	var user *nbum.User

	tI := self.Ctx.Input.Session(nbum.SessionKeyWXResAccessToken)
	resAccessToken, ok := tI.(oauth.ResAccessToken)
	if !ok {
		// 没有 ResAccessToken，需要引导跳转重新登录
		oauth := nbum.DefaultWX.Wechat.GetOauth()
		redirectURL := ""
		var err error
		if self.Ctx.Input.IsAjax() || self.IsSignReq() {
			redirectURL = self.genLoginCallbackURL(beego.AppConfig.String("domain"),
				nbum.AuthLoginWXCallbackPlaceholdOnSucc, nbum.AuthLoginWXCallbackPlaceholdOnSucc)
			redirectURL, err = oauth.GetRedirectURL(redirectURL, "snsapi_userinfo", "STATE")

		} else {
			site := self.Ctx.Input.Site()
			if self.Ctx.Input.Port() != 80 {
				site += ":" + strconv.Itoa(self.Ctx.Input.Port())
			}
			site += self.Ctx.Input.URL()
			redirectURL = self.genLoginCallbackURL(beego.AppConfig.String("domain"), site, "")
			redirectURL, err = oauth.GetRedirectURL(redirectURL, "snsapi_userinfo", "STATE")
		}

		if err != nil {
			dora.Fatal().Err(err).Str("redirectURL", redirectURL).Msg("get wechat redirect url failed")
			self.Ctx.Abort(500, "check login state failed, please contact admin")
		}

		self.Logger().Debug().Bool("IsReqSign", self.IsSignReq()).
			Bool("IsAjax", self.Ctx.Input.IsAjax()).
			Str("redirectURL", redirectURL).Msg("WXLogin will redirect")
		if self.Ctx.Input.IsAjax() || self.IsSignReq() {
			resp := apiresp.NewDetail(nbum.ErrCodeNotLogin, "unauthorized request", redirectURL)
			resp.JSON(self.Ctx)
		} else {
			self.Ctx.Redirect(302, redirectURL)
		}
		panic(beego.ErrAbort)
	} else {
		// 先尝试刷新 access token
		oauth := nbum.DefaultWX.Wechat.GetOauth()
		isOk := false
		var err error

		if resAccessToken, err = oauth.RefreshAccessToken(resAccessToken.RefreshToken); err == nil {
			if user, err = self.userService.GetUserByExternalInfo(
				nbum.UserTypeWX, resAccessToken.OpenID, ""); err == nil {
				self.Ctx.Output.Session(nbum.SessionKeyUser, user)
				isOk = true
			}
		}

		if !isOk {
			self.clearSession()
			dora.Error().Err(err).Interface("oldAccessToken", resAccessToken).Msg("refresh wx oauth access token failed")
			return self.wxLogin()
		}
	}

	return user
}

func (self *LoginStateService) genLoginCallbackURL(domain, succCallback, errorCallback string) string {
	p := url.Values{}
	p.Set(nbum.AuthLoginWXCallbackSuccKey, succCallback)
	if errorCallback != "" {
		p.Set(nbum.AuthLoginWXCallbackErrKey, errorCallback)
	}
	url := domain + nbum.AuthLoginWXCallbackURI + "?" + p.Encode()
	return url
}

func (self *LoginStateService) clearSession() {
	self.Ctx.Input.CruSession.Delete(nbum.SessionKeyUser)
	self.Ctx.Input.CruSession.Delete(nbum.SessionKeyWXResAccessToken)
}
