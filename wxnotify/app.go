package wxnotify

import (
	"encoding/xml"
	"github.com/zohu/zwx"
	"github.com/zohu/zwx/wxcpt"
)

type Context struct {
	*zwx.Context
}

func App(appid string) (*Context, error) {
	c, err := zwx.LoadApp(appid)
	if err != nil {
		return nil, err
	}
	if !c.IsWxMpServe() && !c.IsWxMpSubscribe() && !c.IsWork() && !c.IsWxMiniProgram() {
		return nil, c.Error("", "推送消息仅支持服务号、订阅号、企业号、小程序")
	}
	return &Context{Context: c}, nil
}

func (c *Context) DecodeMessage(p *ReqNotify, recv *wxcpt.BizMsgRecv) (*Message, error) {
	cpt := wxcpt.NewBizMsgCrypt(c.NotifyToken(), c.NotifyEncodingAesKey(), c.AppidMain())
	if cptByte, err := cpt.DecryptMsg(p.MsgSignature, p.Timestamp, p.Nonce, recv); err != nil {
		return nil, err
	} else {
		msg := new(Message)
		msg.Nonce = p.Nonce
		msg.ctx = c
		if err = xml.Unmarshal(cptByte, msg); err != nil {
			return nil, err
		}
		return msg, nil
	}
}
