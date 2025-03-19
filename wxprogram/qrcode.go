package wxprogram

import "github.com/zohu/zwx"

type LineColor struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}
type ReqGetQRCode struct {
	Path       string    `json:"path"`
	Width      int       `json:"width"`
	AutoColor  bool      `json:"auto_color"`
	LineColor  LineColor `json:"line_color"`
	IsHyaline  bool      `json:"is_hyaline"`
	EnvVersion string    `json:"env_version"`
}
type RespGetQRCode struct {
	zwx.WxResponse
	Buffer []byte `json:"buffer"`
}

// GetQRCode
// @Description: 获取小程序码，有数量限制
// @receiver c
// @param req
// @return *RespGetQRCode
// @return error
func (c *Context) GetQRCode(req *ReqGetQRCode) (*RespGetQRCode, error) {
	var resp RespGetQRCode
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiWxa.WithPath("getwxacode")).
		SetAccessToken(c.AccessToken()).
		SetJson(req).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("get_qrcode", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.GetQRCode(req)
		}
		return nil, c.Error("get_qrcode", resp.Errmsg)
	}
	return &resp, nil
}

type ReqGetUnlimitedQRCode struct {
	Scene      string    `json:"scene"`
	Page       string    `json:"page"`
	CheckPath  bool      `json:"check_path"`
	EnvVersion string    `json:"env_version"`
	Width      int       `json:"width"`
	AutoColor  bool      `json:"auto_color"`
	LineColor  LineColor `json:"line_color"`
	IsHyaline  bool      `json:"is_hyaline"`
}

// GetUnlimitedQRCode
// @Description: 获取不限制的小程序码
// @receiver c
// @param req
// @return *RespGetQRCode
// @return error
func (c *Context) GetUnlimitedQRCode(req *ReqGetUnlimitedQRCode) (*RespGetQRCode, error) {
	var resp RespGetQRCode
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiWxa.WithPath("getwxacodeunlimit")).
		SetAccessToken(c.AccessToken()).
		SetJson(req).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("get_limited_qrcode", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.GetUnlimitedQRCode(req)
		}
		return nil, c.Error("get_limited_qrcode", resp.Errmsg)
	}
	return &resp, nil
}

type ReqCreateQRCode struct {
	Path  string `json:"path"`
	Width int    `json:"width"`
}

// CreateQRCode
// @Description: 获取小程序二维码，有数量限制
// @receiver c
// @param req
// @return *RespGetQRCode
// @return error
func (c *Context) CreateQRCode(req *ReqCreateQRCode) (*RespGetQRCode, error) {
	var resp RespGetQRCode
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiCgiBin.WithPath("wxaapp/createwxaqrcode")).
		SetAccessToken(c.AccessToken()).
		SetJson(req).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("create_qrcode", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.CreateQRCode(req)
		}
		return nil, c.Error("create_qrcode", resp.Errmsg)
	}
	return &resp, nil
}
