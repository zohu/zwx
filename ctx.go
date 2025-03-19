package zwx

import (
	"fmt"
	"github.com/zohu/zwx/utils"
	"sync"
	"time"
)

type AppType string

func (a AppType) String() string {
	return string(a)
}

const (
	TypeWxUnknown     AppType = "0"
	TypeWxMpServe     AppType = "1"  // 服务号
	TypeWxMpSubscribe AppType = "2"  // 订阅号
	TypeWxWork        AppType = "3"  // 企业微信
	TypeWxApp         AppType = "4"  // 微信APP生态
	TypeWxMiniApp     AppType = "5"  // 微信小程序
	TypeWxMiniGame    AppType = "6"  // 微信小游戏
	TypeWxOpen        AppType = "7"  // 微信开放平台
	TypeWxVideo       AppType = "8"  // 微信视频号
	TypeWxStore       AppType = "9"  // 微信小店
	TypeWxPay         AppType = "10" // 微信支付
)

type App struct {
	// 应用类型，参考AppType
	AppType AppType `json:"app_type" validate:"required"`
	// 公众号、小程序、小游戏的appid，企业微信的corpid，微信支付mchid
	Appid string `json:"appid" validate:"required"`
	// 公众号、小程序、小游戏、企业微信的secret，微信支付的mch_api_key
	AppSecret string `json:"app_secret" validate:"required"`
	// 订阅号关联的服务号、小程序关联的公众号、企业微信应用关联的企业微信，微信支付关联的业务appid
	MainAppid string `json:"main_appid"`
	// 公众号消息相关的token，微信支付的证书序列号
	Token string `json:"token" validate:"required_if=AppType 10"`
	// 公众号消息相关的EncodingAesKey，微信支付的私钥证书文本内容
	EncodingAesKey string `json:"encoding_aes_key" validate:"required_if=AppType 10"`
	// 微信支付回调地址
	NotifyUri string `json:"notify_uri" validate:"required_if=AppType 10"`
	// DO NOT EDIT, 内部维护字段
	AccessToken string `json:"access_token"`
	// DO NOT EDIT, 内部维护字段
	JsTicket string `json:"js_ticket"`
	// DO NOT EDIT, 内部维护字段
	CardTicket string `json:"card_ticket"`
	// DO NOT EDIT, 内部维护字段
	ExpireTime time.Time `json:"expire_time"`
	// DO NOT EDIT, 内部维护字段
	Retry string `json:"retry"`
}

type Context struct {
	*Wx
	app *App
	sync.Mutex
}

func (c *Context) Appid() string {
	return c.app.Appid
}
func (c *Context) AppidMain() string {
	if c.app.MainAppid != "" {
		return c.app.MainAppid
	}
	return c.app.Appid
}
func (c *Context) AppSecret() string {
	return c.app.AppSecret
}
func (c *Context) IsWxMpServe() bool {
	return c.app.AppType == TypeWxMpServe
}
func (c *Context) IsWxMpSubscribe() bool {
	return c.app.AppType == TypeWxMpSubscribe
}
func (c *Context) IsWork() bool {
	return c.app.AppType == TypeWxWork
}
func (c *Context) IsWxApp() bool {
	return c.app.AppType == TypeWxApp
}
func (c *Context) IsWxMiniProgram() bool {
	return c.app.AppType == TypeWxMiniApp
}
func (c *Context) IsWxMiniGame() bool {
	return c.app.AppType == TypeWxMiniGame
}
func (c *Context) IsWxOpen() bool {
	return c.app.AppType == TypeWxOpen
}
func (c *Context) IsWxVideo() bool {
	return c.app.AppType == TypeWxVideo
}
func (c *Context) IsWxStore() bool {
	return c.app.AppType == TypeWxStore
}
func (c *Context) IsWxPay() bool {
	return c.app.AppType == TypeWxPay
}
func (c *Context) IsDebug() bool {
	return c.debug
}
func (c *Context) Logger() Logger {
	return c.logger
}
func (c *Context) NotifyToken() string {
	return c.app.Token
}
func (c *Context) NotifyEncodingAesKey() string {
	return c.app.EncodingAesKey
}
func (c *Context) NotifyUri() string {
	return c.app.NotifyUri
}
func (c *Context) AccessToken() string {
	if c.app.ExpireTime.Before(time.Now()) {
		c.app.AccessToken = ""
		c.app.JsTicket = ""
		c.app.CardTicket = ""
	}
	if c.app.AccessToken == "" {
		c.NewAccessToken()
	}
	return c.app.AccessToken
}
func (c *Context) JsTicket() string {
	return c.app.JsTicket
}
func (c *Context) CardTicket() string {
	return c.app.CardTicket
}
func (c *Context) NewAccessToken() {
	c.Lock()
	defer c.Unlock()
	if c.app.ExpireTime.Before(time.Now()) {
		c.app.AccessToken = ""
		c.app.JsTicket = ""
		c.app.CardTicket = ""
	}
	switch c.app.AppType {
	case TypeWxMpServe:
		c.newMpToken()
		c.newMpTicket(TicketTypeJs)
		c.newMpTicket(TicketTypeCard)
	case TypeWxMpSubscribe:
		c.newMpToken()
		c.newMpTicket(TicketTypeJs)
	case TypeWxWork:
		c.newWorkToken()
		c.newWorkTicket()
	case TypeWxApp:
		c.newMpToken()
	case TypeWxMiniApp:
		c.newMpToken()
	case TypeWxMiniGame:
		break
	case TypeWxOpen:
		break
	case TypeWxVideo:
		break
	case TypeWxStore:
		break
	case TypeWxPay:
		break
	default:
		c.logger.Errorf("unknown app type: %s", c.app.AppType)
		return
	}
	wxl.Lock()
	defer wxl.Unlock()
	if c.app.AccessToken == "" {
		c.storage.HIncrBy(PrefixApp.Key(c.Appid()), "retry", 1)
	} else {
		c.storage.HSet(PrefixApp.Key(c.Appid()), utils.StructToMap(c.app))
	}
}

// RetryAccessToken
// @Description: 是否可以刷新token并重试(每个app每2分钟只能重试一次)
// @receiver c
// @param ctx
// @param errcode
// @return bool
func (c *Context) RetryAccessToken(errcode int) bool {
	switch errcode {
	case 40014, 41001, 42001, 42007:
		if c.storage.SetNX(PrefixRetry.Key(c.Appid()), "retrying", time.Minute*2) {
			c.NewAccessToken()
			return true
		}
		return false
	default:
		return false
	}
}

func (c *Context) Error(action, message string) error {
	return fmt.Errorf("[%s] %s failed: %s", c.Appid(), action, message)
}
