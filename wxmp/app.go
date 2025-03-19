package wxmp

import "github.com/zohu/zwx"

type Context struct {
	*zwx.Context
}

func App(appid string) (*Context, error) {
	c, err := zwx.LoadApp(appid)
	if err != nil {
		return nil, err
	}
	if !c.IsWxMpServe() && !c.IsWxMpSubscribe() {
		return nil, c.Error("", "非公众号")
	}
	return &Context{Context: c}, nil
}
