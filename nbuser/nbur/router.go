// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (

_ "github.com/NoneBorder/bee-modules/doralog"
"github.com/NoneBorder/bee-modules/nbuser/nbuc"
"github.com/NoneBorder/bee-modules/nbuser/nbus"
"github.com/astaxie/beego"
"github.com/astaxie/beego/context"

)

func init() {
	beego.AddNamespace(
		beego.NewNamespace("/api",
			beego.NSNamespace("/wx",
				beego.NSRouter("/msgreceive/v1", new(nbuc.WXController), "*:MsgReceive"),
			),

			beego.NSNamespace("/u",
				beego.NSRouter("/username", new(nbuc.UserController), "get:UserName"),
				beego.NSRouter("/register", new(nbuc.UserController), "post:Register"),
				beego.NSRouter("/login", new(nbuc.UserController), "post:Login"),
				beego.NSRouter("/logout", new(nbuc.UserController), "get:Logout"),
				beego.NSRouter("/login/wxcb", new(nbuc.UserController), "get:WXLoginCallback"),
			),

			beego.NSNamespace("/common",
				beego.NSRouter("/newVerifyCode", new(nbuc.CommonController), "post:SendVerifyCode"),
			),
		),
	)
}

// SetAuthApi 登录认证接口，出 signauth 外还有登录的 cookie 认证
func SetAuthApi(authApiPatternList []string) {
	for _, api := range authApiPatternList {
		beego.InsertFilter(api, beego.BeforeRouter, func(ctx *context.Context) {
			ls := nbus.NewLoginStateService(ctx)
			ls.GetUser()
		})
	}
}

// SetSignAuthApiPrefix 设置需要走 sign 认证的接口前缀
func SetSignAuthApi(signAuthApiPatternList []string) {
	for _, api := range signAuthApiPatternList {
		beego.InsertFilter(api, beego.BeforeRouter, func(ctx *context.Context) {
			ctx.Request.Header.Set("Nb-Sign", "md5sign")
			ls := nbus.NewLoginStateService(ctx)
			ls.GetUser()
		})
	}
}
