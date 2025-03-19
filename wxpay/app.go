package wxpay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/zohu/zwx"
)

type Context struct {
	*zwx.Context
	sdk *core.Client
}

func App(appid string) (*Context, error) {
	c, err := zwx.LoadApp(appid)
	if err != nil {
		return nil, err
	}
	if !c.IsWxPay() {
		return nil, c.Error("", "此应用非微信商户")
	}
	mchPrivateKey, err := utils.LoadPrivateKey(c.NotifyEncodingAesKey())
	if err != nil {
		return nil, err
	}
	ops := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(c.Appid(), c.NotifyToken(), mchPrivateKey, c.AppSecret()),
	}
	client, err := core.NewClient(context.Background(), ops...)
	if err != nil {
		return nil, err
	}
	return &Context{Context: c, sdk: client}, nil
}

func (c *Context) WxClient() *core.Client {
	return c.sdk
}
func (c *Context) JsapiClient() *jsapi.JsapiApiService {
	return &jsapi.JsapiApiService{Client: c.sdk}
}
func (c *Context) H5Client() *h5.H5ApiService {
	return &h5.H5ApiService{Client: c.sdk}
}
func (c *Context) NativeClient() *native.NativeApiService {
	return &native.NativeApiService{Client: c.sdk}
}
func (c *Context) AppClient() *app.AppApiService {
	return &app.AppApiService{Client: c.sdk}
}
