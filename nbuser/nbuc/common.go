package nbuc

import (
	"encoding/json"

	"github.com/NoneBorder/bee-modules/nbbase"
	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbus"
	"github.com/NoneBorder/dora/apiresp"
)

type CommonController struct {
	BaseController
}

type SendVerifyCodeForm struct {
	Type string
	To   string
}

func (self *CommonController) SendVerifyCode() {
	self.Logger().Debug().Bytes("requestBody", self.Ctx.Input.RequestBody).Msg("request body")

	form := SendVerifyCodeForm{}
	if err := json.Unmarshal(self.Ctx.Input.RequestBody, &form); err != nil {
		self.Logger().Error().Err(err).Msg("parse request body failed")
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "parse request body failed").JSON(self.Ctx)
	}
	self.Logger().Debug().Interface("form", form).Msg("request form")

	if form.Type == "" || form.To == "" {
		apiresp.NewDetail(nbum.ErrCodeBadRequest, "invalid request").JSON(self.Ctx)
	}

	commonService := nbus.NewCommonService(self.Ctx)
	if err := commonService.CreateNewVerifyCode(form.Type, form.To); err != nil {
		self.Logger().Error().Err(err).Msg("create new verify code failed")
		apiresp.NewDetail(nbum.ErrCodeGeneralError, "create verify code failed:"+err.Error()).JSON(self.Ctx)
	}
	apiresp.NewResp(nil).JSON(self.Ctx)
}
