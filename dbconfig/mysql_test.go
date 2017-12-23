package dbconfig

import (
	"strings"
	"testing"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	testDBUserName = "root"
	testDBPass     = "123456"
	testDBName     = "test"

	testCacheTimeout = 5
)

var cfgHandler config.Configer

func init() {
	orm.RegisterDataBase("default", "mysql",
		testDBUserName+":"+testDBPass+"@tcp(127.0.0.1:3306)/"+testDBName+"?charset=utf8", 30)
	orm.RunSyncdb("default", true, true)
	orm.Debug = true

	var err error
	if cfgHandler, err = config.NewConfig("dbconfig", ""); err != nil {
		panic(err)
	}

	cfgHandler.(*MySQLConfiger).SetCacheTimeout(testCacheTimeout * time.Second)
}

func TestGetSet(t *testing.T) {
	key := "test"
	orival := "test"

	var err error
	var val string

	if err = cfgHandler.Set(key, orival); err != nil {
		t.Fatalf("should be ok, but got error: %s", err.Error())
	}

	if val = cfgHandler.String(key); val != orival {
		t.Fatalf("should be equal '%s', but got val '%s'", orival, val)
	}

	orival = "testUpdate"
	if err = cfgHandler.Set(key, orival); err != nil {
		t.Fatalf("should be ok, but got error: %s", err.Error())
	}

	if val = cfgHandler.String(key); val != orival {
		t.Fatalf("should be equal '%s', but got val '%s'", orival, val)
	}

	time.Sleep((testCacheTimeout + 1) * time.Second)
	if val = cfgHandler.String(key); val != orival {
		t.Fatalf("should be equal '%s', but got val '%s'", orival, val)
	}

}

func TestStrings(t *testing.T) {
	key := "testStrings"
	orival := []string{"123", "456", "789"}

	var val []string
	var err error

	val = cfgHandler.DefaultStrings(key, orival)
	if strings.Join(val, ";") != strings.Join(orival, ";") {
		t.Errorf("should be equal to default, but got '%v'", val)
	}

	newval := strings.Join(orival, ";") + ";abc"
	if err = cfgHandler.Set(key, newval); err != nil {
		t.Fatalf("should be ok, but got error: %s", err.Error())
	}

	val = cfgHandler.DefaultStrings(key, orival)
	if strings.Join(val, ";") != newval {
		t.Errorf("should be equal to '%s', but got '%v'", newval, val)
	}
}

func TestInt(t *testing.T) {
	key := "testInt64"
	orival := int64(666)

	var val int64
	var err error

	val = cfgHandler.DefaultInt64(key, orival)
	if val != orival {
		t.Errorf("should be equal to default, but got '%d'", val)
	}

	if err = cfgHandler.Set(key, "8888"); err != nil {
		t.Fatalf("should be ok, but got error: %s", err.Error())
	}

	val = cfgHandler.DefaultInt64(key, orival)
	if val != 8888 {
		t.Errorf("should be equal to 8888, but got '%v'", val)
	}
}
