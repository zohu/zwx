package wxprogram

import (
	"github.com/zohu/zwx"
	"time"
)

type UploadShippingInfoOrderKey struct {
	OrderNumberType int    `json:"order_number_type"`
	TransactionId   string `json:"transaction_id"`
}
type UploadShippingInfoShippingItem struct {
	ItemDesc string `json:"item_desc"`
}
type UploadShippingInfoPayer struct {
	Openid string `json:"openid"`
}
type ParamUploadShippingInfo struct {
	OrderKey      UploadShippingInfoOrderKey       `json:"order_key"`
	LogisticsType int                              `json:"logistics_type"`
	DeliveryMode  int                              `json:"delivery_mode"`
	ShippingList  []UploadShippingInfoShippingItem `json:"shipping_list"`
	UploadTime    string                           `json:"upload_time"`
	Payer         UploadShippingInfoPayer          `json:"payer"`
}

func (c *Context) UploadShippingInfo(openid, itemName, tid string) error {
	var resp zwx.WxResponse
	if err := zwx.NewHttp(zwx.MethodPost, zwx.ApiWxa.WithPath("sec/order/upload_shipping_info")).
		SetAccessToken(c.AccessToken()).
		SetJson(&ParamUploadShippingInfo{
			OrderKey: UploadShippingInfoOrderKey{
				OrderNumberType: 2,
				TransactionId:   tid,
			},
			LogisticsType: 4,
			DeliveryMode:  1,
			ShippingList:  []UploadShippingInfoShippingItem{{ItemDesc: itemName}},
			UploadTime:    time.Now().Format(time.RFC3339),
			Payer:         UploadShippingInfoPayer{Openid: openid},
		}).
		BindJson(&resp).
		Debug(c.IsDebug(), c.Logger()).
		Do(); err != nil {
		return c.Error("upload_shipping_info", err.Error())
	}
	if resp.Errcode != 0 {
		if c.RetryAccessToken(resp.Errcode) {
			return c.UploadShippingInfo(openid, itemName, tid)
		}
		return c.Error("upload_shipping_info", resp.Errmsg)
	}
	return nil
}
