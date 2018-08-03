package nbum

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/astaxie/beego"
)

type VerifyCode struct {
	Code    string
	Created time.Time
}

const notifyBodyTemplate = `<div>您好：</div><div><br></div>
<div>　　您在个付请求的验证码是：<b>%s</b>，%d分钟内有效。</div>
<div>　　<font color="#999999"><i>个付消息通知，请勿直接回复</i></font></div>
<div><br></div><div>感谢您使用个付服务，祝好！</div><div>
<a href="%s"><font color="#00ccff">个付 — 个人支付管理平台</font></a></div>
<div><includetail><!--<![endif]--></includetail></div>
`

const verifyCodeValidDuration = 15

func (self *VerifyCode) GenerateNew() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	self.Code = fmt.Sprintf("%06v", rnd.Int31n(1000000))
	self.Created = time.Now()
}

func (self *VerifyCode) Verify(code string) bool {
	return self.Code != "" && self.Code == code && time.Now().Sub(self.Created) < verifyCodeValidDuration*time.Minute
}

func (self *VerifyCode) NotifyTitle() string {
	return "【个付】验证码"
}

func (self *VerifyCode) NotifyBody() string {
	return fmt.Sprintf(notifyBodyTemplate, self.Code, verifyCodeValidDuration, beego.AppConfig.String("domain"))
}
