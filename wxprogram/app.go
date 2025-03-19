package wxprogram

import (
	"github.com/zohu/zwx"
)

type Context struct {
	*zwx.Context
}

func App(appid string) (*Context, error) {
	c, err := zwx.LoadApp(appid)
	if err != nil {
		return nil, err
	}
	if !c.IsWxMiniProgram() {
		return nil, c.Error("", "this appid is not a miniprogam")
	}
	return &Context{Context: c}, nil
}
