package nbuser

import (
	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbur"
)

func Init(dbAlias, dbTablePrefix string) {
	nbum.DBInit(dbAlias, dbTablePrefix)
	nbur.RegisterRouter()
}
