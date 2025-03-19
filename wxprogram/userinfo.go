package wxprogram

import (
	"github.com/zohu/zwx"
	"github.com/zohu/zwx/wxcpt"
)

type RespGetPluginOpenPId struct {
	zwx.WxResponse
	Openpid string `json:"openpid"`
}

// GetPluginOpenPId
// @Description: 获取插件用户openpid
// @receiver c
// @param code
// @return *RespGetPluginOpenPId
// @return error
func (c *Context) GetPluginOpenPId(code string) (*RespGetPluginOpenPId, error) {
	var resp RespGetPluginOpenPId
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiWxa.WithPath("plugin/get_open_pid")).
		SetAccessToken(c.AccessToken()).
		SetJson(map[string]string{
			"code": code,
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("get_plugin_open_pid", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.GetPluginOpenPId(code)
		}
		return nil, c.Error("get_plugin_open_pid", resp.Errmsg)
	}
	return &resp, nil
}

type RespCheckEncryptedData struct {
	zwx.WxResponse
	Vaild      bool  `json:"vaild"`
	CreateTime int64 `json:"create_time"`
}

// CheckEncryptedData
// @Description: 检查加密信息
// @receiver c
// @param encrypted
// @return *RespCheckEncryptedData
// @return error
func (c *Context) CheckEncryptedData(encrypted string) (*RespCheckEncryptedData, error) {
	var resp RespCheckEncryptedData
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiWxa.WithPath("business/checkencryptedmsg")).
		SetAccessToken(c.AccessToken()).
		SetJson(map[string]string{
			"encrypt_data": encrypted,
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("check_encrypted_data", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.CheckEncryptedData(encrypted)
		}
		return nil, c.Error("check_encrypted_data", resp.Errmsg)
	}
	return &resp, nil
}

type ReqGetPaidUnionid struct {
	Openid        string `json:"openid"`
	TransactionId string `json:"transaction_id"`
	MchId         string `json:"mch_id"`
	OutTradeNo    string `json:"out_trade_no"`
}
type RespGetPaidUnionid struct {
	zwx.WxResponse
	Unionid string `json:"unionid"`
}

// GetPaidUnionid
// @Description: 支付后获取Unionid
// @receiver c
// @param req
// @return *RespGetPaidUnionid
// @return error
func (c *Context) GetPaidUnionid(req *ReqGetPaidUnionid) (*RespGetPaidUnionid, error) {
	var resp RespGetPaidUnionid
	if err := zwx.NewHttp(zwx.MethodGet, zwx.ApiWxa.WithPath("getpaidunionid")).
		SetAccessToken(c.AccessToken()).
		SetQuery(map[string]string{
			"openid":         req.Openid,
			"transaction_id": req.TransactionId,
			"mch_id":         req.MchId,
			"out_trade_no":   req.OutTradeNo,
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("get_paid_unionid", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.GetPaidUnionid(req)
		}
		return nil, c.Error("get_paid_unionid", resp.Errmsg)
	}
	return &resp, nil
}

type RespGetUserEncryptKey struct {
	zwx.WxResponse
	KeyInfoList []struct {
		EncryptKey string `json:"encrypt_key"`
		Version    int    `json:"version"`
		CreateTime int64  `json:"create_time"`
		ExpireIn   int    `json:"expire_in"`
		Iv         string `json:"iv"`
	} `json:"key_info_list"`
}

// GetUserEncryptKey
// @Description: 获取用户encryptKey
// @receiver c
// @param openid
// @param sessionKey
// @return *RespGetUserEncryptKey
// @return error
func (c *Context) GetUserEncryptKey(openid, sessionKey string) (*RespGetUserEncryptKey, error) {
	var resp RespGetUserEncryptKey
	if err := zwx.NewHttp(zwx.MethodGet, zwx.ApiWxa.WithPath("getuserencryptkey")).
		SetAccessToken(c.AccessToken()).
		SetQuery(map[string]string{
			"openid":     openid,
			"signature":  wxcpt.HmacSha256ToBase64("", sessionKey),
			"sig_method": "hmac_sha256",
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("getuserencryptkey", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.GetUserEncryptKey(openid, sessionKey)
		}
		return nil, c.Error("getuserencryptkey", resp.Errmsg)
	}
	return &resp, nil
}

type RespGetPhoneNumber struct {
	zwx.WxResponse
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int64  `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

// GetPhoneNumber
// @Description: 获取手机号
// @receiver c
// @param code
// @param openid
// @return *RespGetPhoneNumber
// @return error
func (c *Context) GetPhoneNumber(code, openid string) (*RespGetPhoneNumber, error) {
	var resp RespGetPhoneNumber
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiWxa.WithPath("business/getuserphonenumber")).
		SetAccessToken(c.AccessToken()).
		SetJson(map[string]string{
			"code":   code,
			"openid": openid,
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return nil, c.Error("get_phone_number", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.GetPhoneNumber(code, openid)
		}
		return nil, c.Error("get_phone_number", resp.Errmsg)
	}
	return &resp, nil
}
