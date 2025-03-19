package wxmp

import "github.com/zohu/zwx"

/**
自定义菜单
*/

type MenuType string // 按钮类型
type MenuKey string  // 按钮功能

const (
	MenuTypeView        MenuType = "view"        // 视图按钮
	MenuTypeClick       MenuType = "click"       // 点击按钮
	MenuTypeMiniprogram MenuType = "miniprogram" // 小程序按钮
)
const (
	MenuKeyClick              MenuKey = "click"
	MenuKeyView               MenuKey = "view"
	MenuKeyScancodePush       MenuKey = "scancode_push"        // 扫码推事件
	MenuKeyScancodeWaitmsg    MenuKey = "scancode_waitmsg"     // 扫码推事件且弹出“消息接收中”提示框
	MenuKeyPicSysphoto        MenuKey = "pic_sysphoto"         // 弹出系统拍照发图
	MenuKeyPicPhotoOrAlbum    MenuKey = "pic_photo_or_album"   // 弹出拍照或者相册发图
	MenuKeyPicWeixin          MenuKey = "pic_weixin"           // 弹出微信相册发图器
	MenuKeyLocationSelect     MenuKey = "location_select"      // 弹出地理位置选择器
	MenuKeyMediaId            MenuKey = "media_id"             // 下发消息（除文本消息）
	MenuKeyArticleId          MenuKey = "article_id"           // 微信客户端将会以卡片形式，下发开发者在按钮中填写的图文消息
	MenuKeyArticleViewLimited MenuKey = "article_view_limited" // 类似 view_limited，但不使用 media_id 而使用 article_id
)

type MenuButtonItem struct {
	Type      MenuType         `json:"type,omitempty"`
	Name      string           `json:"name"`
	Key       MenuKey          `json:"key,omitempty"`
	Url       string           `json:"url,omitempty"`
	Appid     string           `json:"appid,omitempty"`
	MediaId   string           `json:"media_id,omitempty"`
	ArticleId string           `json:"article_id,omitempty"`
	Pagepath  string           `json:"pagepath,omitempty"`
	SubButton []MenuButtonItem `json:"sub_button,omitempty"`
}
type Menu struct {
	Button []MenuButtonItem `json:"button"`
}

// MenuDiyMatchRule
// @Description: 性别、国家、省市区、语言 官方已经废除，不再支持
type MenuDiyMatchRule struct {
	TagId              string `json:"tag_id,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
}

type MenuDiy struct {
	Button    []MenuButtonItem `json:"button"`
	Matchrule MenuDiyMatchRule `json:"matchrule"`
}

func (c *Context) MenuAdd(menu *Menu) error {
	var resp zwx.WxResponse
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiCgiBin.WithPath("menu/create")).
		SetAccessToken(c.AccessToken()).
		SetJson(menu).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return c.Error("menu add", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.MenuAdd(menu)
		}
		return c.Error("menu add", resp.Errmsg)
	}
	return nil
}
