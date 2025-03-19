package zwx

import "time"

type ResAccessToken struct {
	WxResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type TicketType string

const (
	TicketTypeJs   TicketType = "jsapi"
	TicketTypeCard TicketType = "wx_card"
)

type ResTicket struct {
	WxResponse
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// -------------------------------mp-------------------------------

func (c *Context) newMpToken() {
	var resp ResAccessToken
	if err := NewHttp(MethodGet, ApiCgiBin.WithPath("token")).
		SetQuery(map[string]string{
			"grant_type": "client_credential",
			"appid":      c.AppidMain(),
			"secret":     c.AppSecret(),
		}).
		BindJson(&resp).
		Debug(c.debug, c.logger).
		Do(); err != nil {
		c.logger.Errorf("%s request access_token failed：%s", c.AppidMain(), err.Error())
		return
	}
	if resp.Errcode != 0 {
		c.logger.Errorf("%s request access_token failed：%s", c.AppidMain(), resp.Errmsg)
		return
	}
	c.app.AccessToken = resp.AccessToken
	c.app.ExpireTime = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
}
func (c *Context) newMpTicket(t TicketType) {
	if c.app.AccessToken == "" {
		return
	}
	var resp ResTicket
	if err := NewHttp(MethodGet, ApiCgiBin.WithPath("ticket/getticket")).
		SetAccessToken(c.app.AccessToken).
		SetQuery(map[string]string{
			"type": string(t),
		}).
		BindJson(&resp).
		Debug(c.debug, c.logger).
		Do(); err != nil {
		c.logger.Errorf("%s request ticket failed：%s", c.app.AccessToken, err.Error())
		return
	}
	if resp.Errcode != 0 {
		c.logger.Errorf("%s request ticket failed：%s", c.app.AccessToken, resp.Errmsg)
		return
	}
	switch t {
	case TicketTypeJs:
		c.app.JsTicket = resp.Ticket
	case TicketTypeCard:
		c.app.CardTicket = resp.Ticket
	}
}

// -------------------------------work-------------------------------

func (c *Context) newWorkToken() {
	var resp ResAccessToken
	if err := NewHttp(MethodGet, ApiWorkCgiBin.WithPath("gettoken")).
		SetQuery(map[string]string{
			"corpid":     c.AppidMain(),
			"corpsecret": c.AppSecret(),
		}).
		BindJson(&resp).
		Debug(c.debug, c.logger).
		Do(); err != nil {
		c.logger.Errorf("%s request work access_token failed：%s", c.AppidMain(), err.Error())
		return
	}
	if resp.Errcode != 0 {
		c.logger.Errorf("%s request work access_token failed：%s", c.AppidMain(), resp.Errmsg)
		return
	}
	c.app.AccessToken = resp.AccessToken
	c.app.ExpireTime = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
}
func (c *Context) newWorkTicket() {
	if c.app.AccessToken == "" {
		return
	}
	var resp ResTicket
	if err := NewHttp(MethodGet, ApiWorkCgiBin.WithPath("ticket/get")).
		SetAccessToken(c.app.AccessToken).
		SetQuery(map[string]string{
			"type": "agent_config",
		}).
		BindJson(&resp).
		Debug(c.debug, c.logger).
		Do(); err != nil {
		c.logger.Errorf("%s request ticket failed：%s", c.app.AccessToken, err.Error())
		return
	}
	if resp.Errcode != 0 {
		c.logger.Errorf("%s request ticket failed：%s", c.app.AccessToken, resp.Errmsg)
		return
	}
	c.app.JsTicket = resp.Ticket
}
