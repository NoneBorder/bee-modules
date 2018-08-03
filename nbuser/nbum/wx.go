package nbum

import (
	"github.com/astaxie/beego"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
)

/*

WX 微信 lib

## wechat config example
wxRedisHost = 127.0.0.1:6379
wxAppID = wx62e37b113489e3d4
wxAppSecret = 32075673ab541e977ba3176f56fcdac5
wxToken = w77Y96KxdQ69AU8sVN97BCXLvRwx7sjt
wxEncodingAESKey = Nybu6rYRWFVgAptXeKAh3TX3JvgFMheZnexxN68kxtsxamTMsvF7GDE7E89FPmG8

 */
type WX struct {
	config *wechat.Config
	Wechat *wechat.Wechat
}

var DefaultWX *WX

func InitWechatConfig() {
	DefaultWX = new(WX)
	DefaultWX.Init()
}

func (self *WX) Init() {
	redisCache := cache.NewRedis(
		&cache.RedisOpts{Host: beego.AppConfig.DefaultString("wxRedisHost", "127.0.0.1:6379")},
	)

	self.config = &wechat.Config{
		AppID:          beego.AppConfig.String("wxAppID"),
		AppSecret:      beego.AppConfig.String("wxAppSecret"),
		Token:          beego.AppConfig.String("wxToken"),
		EncodingAESKey: beego.AppConfig.String("wxEncodingAESKey"),
		Cache:          redisCache,
	}
	self.Wechat = wechat.NewWechat(self.config)
}
