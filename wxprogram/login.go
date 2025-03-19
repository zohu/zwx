package wxprogram

import (
	"github.com/zohu/zwx"
	"github.com/zohu/zwx/wxcpt"
)

type ResCode2Session struct {
	zwx.WxResponse
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

// Code2Session
// @Description: 小程序登录
// @receiver c
// @param code
// @return *ResCode2Session
// @return error
func (c *Context) Code2Session(code string) (*ResCode2Session, error) {
	var resp ResCode2Session
	if err := zwx.NewHttp(zwx.MethodGet, zwx.ApiSns.WithPath("jscode2session")).
		SetQuery(map[string]string{
			"appid":      c.Appid(),
			"secret":     c.AppSecret(),
			"js_code":    code,
			"grant_type": "authorization_code",
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("code2session", err.Error())
	}
	if resp.Errcode != 0 {
		return nil, c.Error("code2session", resp.Errmsg)
	}
	return &resp, nil
}

// CheckSessionKey
// @Description: 检验登录态
// @receiver c
// @param openid
// @param sessionKey
// @return *zwx.WxResponse
// @return error
func (c *Context) CheckSessionKey(openid, sessionKey string) (*zwx.WxResponse, error) {
	var resp zwx.WxResponse
	if err := zwx.NewHttp(zwx.MethodGet, zwx.ApiWxa.WithPath("checksession")).
		SetAccessToken(c.AccessToken()).
		SetQuery(map[string]string{
			"openid":     openid,
			"signature":  wxcpt.HmacSha256ToBase64("", sessionKey),
			"sig_method": "hmac_sha256",
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("checksession", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.CheckSessionKey(openid, sessionKey)
		}
		return nil, c.Error("checksession", resp.Errmsg)
	}
	return &resp, nil
}

type RespResetUserSessionKey struct {
	zwx.WxResponse
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

// ResetUserSessionKey
// @Description: 重置登录态
// @receiver c
// @param openid
// @param sessionKey
// @return *RespResetUserSessionKey
// @return error
func (c *Context) ResetUserSessionKey(openid, sessionKey string) (*RespResetUserSessionKey, error) {
	var resp RespResetUserSessionKey
	if err := zwx.NewHttp(zwx.MethodGet, zwx.ApiWxa.WithPath("resetusersessionkey")).
		SetAccessToken(c.AccessToken()).
		SetQuery(map[string]string{
			"openid":     openid,
			"signature":  wxcpt.HmacSha256ToBase64("", sessionKey),
			"sig_method": "hmac_sha256",
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("reset checksession", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.ResetUserSessionKey(openid, sessionKey)
		}
		return nil, c.Error("reset checksession", resp.Errmsg)
	}
	return &resp, nil
}
