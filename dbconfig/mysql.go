package dbconfig

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
)

// ConfigTable 是配置的 DB 表结构，主要的配置都存储在 DB 中
type ConfigTable struct {
	Id      int
	Key     string    `orm:"unique"`
	Value   string    `orm:"type(text)"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Comment string    // 配置说明
}

const (
	// cacheKeyPrefix 用户配置缓存 key 的前缀
	cacheKeyPrefix = "dbconfig:"

	// defaultUseDB  设置使用的 DB， 默认是 default
	defaultUseDB = "default"

	// defaultCacheTimeout  默认缓存时间是 60s
	defaultCacheTimeoutSecond = 60
)

var (
	// Err
	ErrNotSupport error = errors.New("not support for dbconfig")
)

func init() {
	// 注册配置模型到 DB
	orm.RegisterModel(new(ConfigTable))

	// 注册到 beego config
	config.Register("dbconfig", new(MySQLConfig))
}

// MySQLConfig 是一个以 MySQL 作为存储后端动态配置组件，完全继承了 beego Config 接口
type MySQLConfig struct{}

func (self *MySQLConfig) Parse(filename string) (config.Configer, error) {
	return self.ParseData([]byte(""))
}

func (self *MySQLConfig) ParseData(data []byte) (config.Configer, error) {
	cacheHandler, err := cache.NewCache("memory", `{"interval":60}`)

	cfg := &MySQLConfiger{
		useDB:        defaultUseDB,
		cacheTimeout: defaultCacheTimeoutSecond * time.Second,
		cacheHandler: cacheHandler,
	}

	return cfg, err
}

// MySQLConfiger 是 MySQL 数据读取的实例，只有在读取数据时使用
type MySQLConfiger struct {
	useDB        string
	cacheHandler cache.Cache
	cacheTimeout time.Duration
}

// SetUseDB 设置要使用的 DB，该 DB 必须已经在主程序中初始化
func (self *MySQLConfiger) SetUseDB(db string) {
	self.useDB = db
}

// SetCacheTimeout 设置多次读取同一个 key 的缓存时间
func (self *MySQLConfiger) SetCacheTimeout(timeout time.Duration) {
	self.cacheTimeout = timeout
}

func (self *MySQLConfiger) Set(key, val string) error {
	cfg := ConfigTable{
		Key:   key,
		Value: val,
	}

	// 保存到 DB
	o := orm.NewOrm()
	o.Using(self.useDB)
	if created, _, err := o.ReadOrCreate(&cfg, "Key"); err == nil {
		if !created {
			// 没有新插入，需要更新
			cfg.Value = val // 重新设置 val， 保证正确
			if _, err := o.Update(&cfg, "Value", "Updated"); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	// 保存到 cache
	cacheKey := cacheKeyPrefix + key
	self.cacheHandler.Put(cacheKey, cfg, self.cacheTimeout) // 忽略缓存错误

	return nil
}

func (self *MySQLConfiger) String(key string) string {
	cacheKey := cacheKeyPrefix + key
	cfg := ConfigTable{Key: key}

	if valObj := self.cacheHandler.Get(cacheKey); valObj == nil {
		// cache 没有数据，从 DB 读取
		o := orm.NewOrm()
		o.Using(self.useDB)
		if err := o.Read(&cfg, "Key"); err == nil {
			// 写入缓存
			self.cacheHandler.Put(cacheKey, cfg, self.cacheTimeout)
		}
	} else {
		// cache 有数据
		cfg, _ = valObj.(ConfigTable)
	}

	return cfg.Value
}

func (self *MySQLConfiger) DefaultString(key, defaultVal string) string {
	if val := self.String(key); val != "" {
		return val
	}
	return defaultVal
}

func (self *MySQLConfiger) Strings(key string) []string {
	stringVal := self.String(key)
	if stringVal == "" {
		return nil
	}
	return strings.Split(stringVal, ";")
}

func (self *MySQLConfiger) DefaultStrings(key string, defaultVal []string) []string {
	if val := self.Strings(key); val != nil {
		return val
	}
	return defaultVal
}

func (self *MySQLConfiger) Int(key string) (int, error) {
	val := self.String(key)
	if val == "" {
		return 0, errors.New("not exist key:" + key)
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New("not int value")
	}

	return v, nil
}

func (self *MySQLConfiger) DefaultInt(key string, defaultVal int) int {
	if val, err := self.Int(key); err == nil {
		return val
	}
	return defaultVal
}

func (self *MySQLConfiger) Int64(key string) (int64, error) {
	val := self.String(key)
	if val == "" {
		return 0, errors.New("not exist key:" + key)
	}

	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.New("not int64 value")
	}

	return v, nil
}

func (self *MySQLConfiger) DefaultInt64(key string, defaultVal int64) int64 {
	if val, err := self.Int64(key); err == nil {
		return val
	}
	return defaultVal
}

func (self *MySQLConfiger) Bool(key string) (bool, error) {
	val := self.String(key)
	switch strings.ToLower(val) {
	case "1", "t", "true", "yes", "y", "on":
		return true, nil
	case "0", "f", "false", "no", "n", "off":
		return false, nil
	case "":
		return false, errors.New("not exist key:" + key)
	default:
		return false, errors.New("not bool value: " + val)
	}
}

func (self *MySQLConfiger) DefaultBool(key string, defaultVal bool) bool {
	if val, err := self.Bool(key); err == nil {
		return val
	}
	return defaultVal
}

func (self *MySQLConfiger) Float(key string) (float64, error) {
	val := self.String(key)
	if val == "" {
		return 0, errors.New("not exist key:" + key)
	}

	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errors.New("not float64 value")
	}

	return v, nil
}

func (self *MySQLConfiger) DefaultFloat(key string, defaultVal float64) float64 {
	if val, err := self.Float(key); err == nil {
		return val
	}
	return defaultVal
}

func (self *MySQLConfiger) DIY(key string) (interface{}, error) {
	return nil, ErrNotSupport
}

func (self *MySQLConfiger) GetSection(section string) (map[string]string, error) {
	return nil, ErrNotSupport
}

func (self *MySQLConfiger) SaveConfigFile(filename string) error {
	return ErrNotSupport
}
