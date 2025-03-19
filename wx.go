package zwx

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/zohu/zwx/utils"
	"sync"
	"time"
)

var wxl sync.Mutex

type Wx struct {
	debug              bool
	logger             Logger
	storage            Storage
	accessTokenRefresh time.Duration
}

var wx *Wx

// New
// @Description: 初始化管理器
// @param options
func New(options *Options) {
	options.Validate()
	wx = &Wx{
		debug:              options.Debug,
		logger:             &logger{debug: options.Debug, Logger: options.Logger},
		storage:            &storage{prefix: options.StoragePrefix, s: options.Storage},
		accessTokenRefresh: options.AccessTokenRefresh,
	}
	if options.AlwaysCleanBeforeStart {
		for _, appid := range Appids() {
			DeleteApp(appid)
		}
	} else {
		go wx.refreshAccessTokenMember()
	}
	go wx.refreshAccessToken()
	wx.logger.Infof("init zwx success")
}

func (wx *Wx) refreshAccessToken() {
	defer func() {
		if r := recover(); r != nil {
			wx.logger.Errorf("refresh token panic: %v", r)
			wx.refreshAccessToken()
		}
	}()
	for range time.NewTicker(wx.accessTokenRefresh).C {
		wx.refreshAccessTokenMember()
	}
}
func (wx *Wx) refreshAccessTokenMember() {
	appids := wx.storage.SMembers(PrefixAppList.Key())
	for _, appid := range appids {
		if c, err := LoadApp(appid); err != nil {
			wx.logger.Errorf("load app %s error: %v", appid, err)
		} else {
			wx.logger.Debugf("refresh %s access token", appid)
			c.NewAccessToken()
		}
	}
	wx.logger.Debugf("wx token refreshed")
}

// LoadApp
// @Description: 获取APP实例
// @param ctx
// @param appid
// @return *Context
// @return error
func LoadApp(appid string) (*Context, error) {
	wxl.Lock()
	defer wxl.Unlock()
	if m := wx.storage.HGetAll(PrefixApp.Key(appid)); len(m) == 0 {
		return nil, fmt.Errorf("appid %s not found", appid)
	} else {
		app := new(App)
		d, _ := sonic.Marshal(m)
		_ = sonic.Unmarshal(d, app)
		return &Context{app: app, Wx: wx}, nil
	}
}

// CreateApp
// @Description: 创建并托管APP实例
// @param ctx
// @param app
// @return error
func CreateApp(app App) error {
	if err := utils.Validate(app); err != nil {
		return fmt.Errorf("create app %s error: %v", app.Appid, err)
	}
	wxl.Lock()
	app.Retry = "0"
	app.ExpireTime = time.Now()
	wx.storage.SAdd(PrefixAppList.Key(), app.Appid)
	wx.storage.HSet(PrefixApp.Key(app.Appid), utils.StructToMap(app))
	wxl.Unlock()
	if c, err := LoadApp(app.Appid); err != nil {
		return fmt.Errorf("create app %s error: %v", app.Appid, err)
	} else {
		c.NewAccessToken()
	}
	wx.logger.Debugf("create app %s success", app.Appid)
	return nil
}

// DeleteApp
// @Description: 停止托管APP实例
// @param ctx
// @param appid
func DeleteApp(appid string) {
	wxl.Lock()
	defer wxl.Unlock()
	wx.logger.Debugf("delete app: %s", appid)
	wx.storage.SRem(PrefixAppList.Key(), appid)
	wx.storage.Del(PrefixApp.Key(appid))
}

// Appids
// @Description: 获取已托管APPID列表
// @param ctx
// @return []string
func Appids() []string {
	wxl.Lock()
	defer wxl.Unlock()
	return wx.storage.SMembers(PrefixAppList.Key())
}

// logger
// @Description: 覆写debugf，支持debug模式
type logger struct {
	debug bool
	Logger
}

func (l *logger) Debugf(format string, v ...any) {
	if l.debug {
		l.Logger.Debugf(format, v...)
	}
}

// storage
// @Description: 覆写storage，以支持prefix
type storage struct {
	prefix string
	s      Storage
}

func (s *storage) Get(key string) string {
	return s.s.Get(s.pre(key))
}
func (s *storage) Del(key string) {
	s.s.Del(s.pre(key))
}
func (s *storage) SetEX(key string, val string, expire time.Duration) {
	s.s.SetEX(s.pre(key), val, expire)
}
func (s *storage) SetNX(key string, val string, expire time.Duration) bool {
	return s.s.SetNX(s.pre(key), val, expire)
}
func (s *storage) SAdd(key string, val ...string) {
	s.s.SAdd(s.pre(key), val...)
}
func (s *storage) SRem(key string, val ...string) {
	s.s.SRem(s.pre(key), val...)
}
func (s *storage) SMembers(key string) []string {
	return s.s.SMembers(s.pre(key))
}
func (s *storage) HSet(key string, val map[string]string) {
	s.s.HSet(s.pre(key), val)
}
func (s *storage) HGetAll(key string) map[string]string {
	return s.s.HGetAll(s.pre(key))
}
func (s *storage) HIncrBy(key string, field string, incr int64) {
	s.s.HIncrBy(s.pre(key), field, incr)
}
func (s *storage) pre(key string) string {
	if s.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", s.prefix, key)
}
