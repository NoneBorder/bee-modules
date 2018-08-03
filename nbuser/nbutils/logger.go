package nbutils

import (
	"github.com/NoneBorder/bee-modules/doralog"
	"github.com/NoneBorder/dora"
	"github.com/astaxie/beego/context"
	"github.com/rs/zerolog"
)

type CommonLogFieldsWrapper func(c zerolog.Context) zerolog.Context

func BaseLoggerSetup(ctx *context.Context, commonLogFieldsWrapper CommonLogFieldsWrapper) *zerolog.Logger {
	var logger *zerolog.Logger

	if ctx == nil {
		logger = &dora.Logger
	} else {
		logger = doralog.H(ctx)
	}

	if commonLogFieldsWrapper != nil {
		c := logger.With()
		c = commonLogFieldsWrapper(c)
		_logger := c.Logger()
		logger = &_logger
	}

	return logger
}
