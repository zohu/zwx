package wxpay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/zohu/zwx/utils"
	"time"
)

type PayType int

const (
	PayTypeJs          PayType = iota + 1 // jsapi支付
	PayTypeApp                            // app支付
	PayTypeH5                             // h5支付
	PayTypeNative                         // C扫B支付
	PayTypeMiniProgram                    // 小程序支付
	PayTypeScanQrcode                     // B扫C收款
)

type GoodsDetail struct {
	MerchantGoodsId  *string `json:"merchant_goods_id"`
	WechatpayGoodsId *string `json:"wechatpay_goods_id,omitempty"`
	GoodsName        *string `json:"goods_name,omitempty"`
	Quantity         *int64  `json:"quantity"`
	UnitPrice        *int64  `json:"unit_price"`
}
type Detail struct {
	CostPrice   *int64        `json:"cost_price,omitempty"`
	InvoiceId   *string       `json:"invoice_id,omitempty"`
	GoodsDetail []GoodsDetail `json:"goods_detail,omitempty"`
}
type StoreInfo struct {
	Id       *string `json:"id"`
	Name     *string `json:"name,omitempty"`
	AreaCode *string `json:"area_code,omitempty"`
	Address  *string `json:"address,omitempty"`
}
type SceneInfo struct {
	PayerClientIp *string    `json:"payer_client_ip"`
	DeviceId      *string    `json:"device_id,omitempty"`
	StoreInfo     *StoreInfo `json:"store_info,omitempty"`
}
type ReqPrepay struct {
	// 支付方式
	PayType PayType `json:"pay_type" validate:"required"`
	// 公众号ID
	Appid *string `json:"appid" validate:"required"`
	// 用户openid
	Openid *string `json:"openid" validate:"required"`
	// 商品描述
	Description *string `json:"description" validate:"required,max=127"`
	// 商户订单号
	OutTradeNo *string `json:"out_trade_no" validate:"required"`
	// 金额，分
	Amount *int64 `json:"amount" validate:"min=1"`
	// 订单失效时间
	TimeExpire *time.Time `json:"time_expire,omitempty"`
	// 附加数据
	Attach *string `json:"attach,omitempty"`
	// 商品标记，代金券或立减优惠功能的参数。
	GoodsTag *string `json:"goods_tag,omitempty"`
	// 指定支付方式
	LimitPay []string `json:"limit_pay,omitempty"`
	// 传入true时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效。
	SupportFapiao *bool `json:"support_fapiao,omitempty"`
	// 优惠功能
	Detail *Detail `json:"detail,omitempty"`
	// 场景信息
	SceneInfo *SceneInfo `json:"scene_info,omitempty"`
	// 分账标识
	ProfitSharing *bool `json:"profit_sharing,omitempty"`
}
type RespPrepay struct {
	PrepayId string `json:"prepay_id,omitempty"`
}

func (c *Context) Prepay(req *ReqPrepay) (*RespPrepay, error) {
	if err := utils.Validate(req); err != nil {
		return nil, c.Error("prepay", err.Error())
	}
	ctx := context.Background()
	var resp RespPrepay
	switch req.PayType {
	case PayTypeJs:
		param := jsapi.PrepayRequest{
			Appid:         utils.FirstTruth(req.Appid, utils.Ptr(c.AppidMain())),
			Mchid:         utils.Ptr(c.Appid()),
			Description:   req.Description,
			OutTradeNo:    req.OutTradeNo,
			TimeExpire:    req.TimeExpire,
			Attach:        req.Attach,
			NotifyUrl:     core.String(c.NotifyUri()),
			GoodsTag:      req.GoodsTag,
			LimitPay:      req.LimitPay,
			SupportFapiao: req.SupportFapiao,
			Amount: &jsapi.Amount{
				Total:    req.Amount,
				Currency: utils.Ptr("CNY"),
			},
			Payer:      &jsapi.Payer{Openid: req.Openid},
			Detail:     req.Detail,
			SceneInfo:  req.SceneInfo,
			SettleInfo: req.SettleInfo,
		}
		if req.Detail != nil {
			param.Detail = &jsapi.Detail{
				CostPrice:   req.Detail.CostPrice,
				InvoiceId:   req.Detail.InvoiceId,
				GoodsDetail: req.Detail.GoodsDetail,
			}
		}
		res, _, err := c.JsapiClient().Prepay(ctx, param)
		if err != nil {
			return nil, c.Error("prepay", err.Error())
		}
		resp.PrepayId = *res.PrepayId
	case PayTypeApp:
		c.AppClient().Prepay(ctx, app.PrepayRequest{
			Appid:         utils.FirstTruth(req.Appid, utils.Ptr(c.AppidMain())),
			Mchid:         utils.Ptr(c.Appid()),
			Description:   req.Description,
			OutTradeNo:    req.OutTradeNo,
			TimeExpire:    req.TimeExpire,
			Attach:        req.Attach,
			NotifyUrl:     utils.Ptr(c.NotifyUri()),
			GoodsTag:      req.GoodsTag,
			LimitPay:      req.LimitPay,
			SupportFapiao: req.SupportFapiao,
			Amount: &app.Amount{
				Total:    req.Amount,
				Currency: utils.Ptr("CNY"),
			},
			Detail:     req.Detail,
			SceneInfo:  req.SceneInfo,
			SettleInfo: req.SettleInfo,
		})
	case PayTypeH5:
	case PayTypeNative:
	case PayTypeMiniProgram:
	case PayTypeScanQrcode:
	}
	return &resp, nil
}
