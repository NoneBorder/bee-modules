package doralog

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/NoneBorder/dora"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/rs/zerolog"
)

const (
	LogHandlerKey = "doraLogHandler"
)

// trace id key name in the http header
var LogTraceIDHeaderKeyName string = "doraLogTraceID"

func init() {
	registerFilter()
}

func registerFilter() {
	beego.InsertFilter("/*", beego.BeforeStatic, func(ctx *context.Context) {
		loggerI := ctx.Input.GetData(LogHandlerKey)
		if loggerI != nil {
			// exists logger, return
			return
		}
		traceID := getTraceID(ctx)
		logger := dora.With().Str("traceID", traceID).Logger()
		ctx.Input.SetData(LogHandlerKey, logger)
	})
}

func getTraceID(ctx *context.Context) string {
	traceID := ctx.Input.Header(LogTraceIDHeaderKeyName)
	if traceID == "" {
		// generate when empty
		hashStr := time.Now().Format(time.RFC3339Nano) + ctx.Input.URL() + ctx.Input.UserAgent() +
			ctx.Input.Cookie(beego.BConfig.WebConfig.Session.SessionName)
		hash := md5.Sum([]byte(hashStr))
		traceID = fmt.Sprintf("%X", hash[:])
	}
	ctx.Input.SetData(LogTraceIDHeaderKeyName, traceID)
	ctx.Output.Header(LogTraceIDHeaderKeyName, traceID)
	return traceID
}

func H(ctx *context.Context) *zerolog.Logger {
	loggerI := ctx.Input.GetData(LogHandlerKey)
	logger, ok := loggerI.(zerolog.Logger)
	if loggerI == nil || !ok {
		return &dora.Logger
	}
	return &logger
}
