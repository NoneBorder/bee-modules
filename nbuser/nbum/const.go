package nbum

const (
	SessionKeyUser             = "user"
	SessionKeyWXResAccessToken = "wxResAccessToken"
	SessionKeyVerifyCode       = "verifyCode"
)

const (
	HeaderUserType = "Nb-User-Type"
	HeaderNBSign   = "Nb-Sign"
)

const (
	NBSignTypeMd5Sign    = "md5sign"
	SignParamTS          = "ts"
	SignParamUserId      = "userId"
	SignParamSign        = "sign"
	SignParamRequestBody = "requestbody"
	SignParamSecretKey   = "_secretKey_"

	CtxInputDataKeyUser = "user"
)

const (
	UserTypeWX     = "wechat"
	UserTypeEmail  = "email"
	UserTypeUnkown = "unkown"
)

const (
	AuthLoginWXCallbackURI             = "/api/u/login/wxcb"
	AuthLoginWXCallbackSuccKey         = "succcb"
	AuthLoginWXCallbackPlaceholdOnSucc = "SUCC_REDIRECT"
	AuthLoginWXCallbackErrKey          = "errcb"
	AuthLoginWXCallbackPlaceholdOnErr  = "ERR_REDIRECT"
)

