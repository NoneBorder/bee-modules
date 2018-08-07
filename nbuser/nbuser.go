package nbuser

import (
	"github.com/NoneBorder/bee-modules/nbuser/nbum"
	"github.com/NoneBorder/bee-modules/nbuser/nbur"
	_ "github.com/NoneBorder/bee-modules/nbuser/nbut"
)

func Init(dbAlias, dbTablePrefix string) {
	nbum.DBInit(dbAlias, dbTablePrefix)
	nbur.RegisterRouter()
}
