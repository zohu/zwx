package wxnotify

import (
	"encoding/xml"
)

type ReqNotify struct {
	MsgSignature string `json:"msg_signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Echostr      string `json:"echostr,omitempty"`
	EncryptType  string `json:"encrypt_type,omitempty"`
}

// MessageType
// @Description: 消息类型
type MessageType string  // 消息类型
type MessageEvent string // 事件类型

// 消息类型
const (
	// 基础类型

	MessageTypeText            MessageType = "text"                      // 文本消息
	MessageTypeImage           MessageType = "image"                     // 图片消息
	MessageTypeVoice           MessageType = "voice"                     // 语音消息
	MessageTypeVideo           MessageType = "video"                     // 视频消息
	MessageTypeMusic           MessageType = "music"                     // 音乐消息
	MessageTypeNews            MessageType = "news"                      // 图文消息
	MessageTypeShortvideo      MessageType = "shortvideo"                // 小视频消息
	MessageTypeLocation        MessageType = "location"                  // 地理位置消息
	MessageTypeLink            MessageType = "link"                      // 链接消息
	MessageTypeMiniprogrampage MessageType = "miniprogrampage"           // 小程序卡片消息
	MessageTypeTransfer        MessageType = "transfer_customer_service" // 消息消息转发到客服
	MessageTypeEvent           MessageType = "event"                     // 事件消息

	// 企业号

	MessageTypeUpdateButton       MessageType = "update_button"        // 更新点击用户的按钮文案
	MessageTypeUpdateTemplateCard MessageType = "update_template_card" // 企业微信被动回复
)

// 事件类型
const (

	// 公众号事件

	MessageEventSubscribe               MessageEvent = "subscribe"                  // 关注事件
	MessageEventUnsubscribe             MessageEvent = "unsubscribe"                // 取关事件
	MessageEventScan                    MessageEvent = "SCAN"                       // 扫描二维码事件
	MessageEventLocation                MessageEvent = "LOCATION"                   // 上报地理位置事件
	MessageEventClick                   MessageEvent = "CLICK"                      // 点击菜单拉取消息时的事件
	MessageEventView                    MessageEvent = "VIEW"                       // 点击菜单跳转链接时的事件推送
	MessageEventScancodePush            MessageEvent = "scancode_push"              // 扫码推事件的事件推送
	MessageEventScancodeWaitmsg         MessageEvent = "scancode_waitmsg"           // 扫码推事件且弹出“消息接收中”提示框的事件推送
	MessageEventPicSysphoto             MessageEvent = "pic_sysphoto"               // 弹出系统拍照发图的事件推送
	MessageEventPicPhotoOrAlbum         MessageEvent = "pic_photo_or_album"         // 弹出拍照或者相册发图的事件推送
	MessageEventPicWeixin               MessageEvent = "pic_weixin"                 // 弹出微信相册发图器的事件推送
	MessageEventLocationSelect          MessageEvent = "location_select"            // 弹出地理位置选择器的事件推送
	MessageEventTemplateSendJobFinish   MessageEvent = "TEMPLATESENDJOBFINISH"      // 发送模板消息推送通知
	MessageEventMassSendJobFinish       MessageEvent = "MASSSENDJOBFINISH"          // 群发消息推送通知
	MessageEventWxaMediaCheck           MessageEvent = "wxa_media_check"            // 异步校验图片/音频是否含有违法违规内容推送事件
	MessageEventSubscribeMsgPopupEvent  MessageEvent = "subscribe_msg_popup_event"  // 订阅通知事件推送
	MessageEventSubscribeMsgChangeEvent MessageEvent = "subscribe_msg_change_event" // 订阅通知用户管理
	MessageEventSubscribeMsgSentEvent   MessageEvent = "subscribe_msg_sent_event"   // 订阅通知发送订阅通知
	MessageEventPublishJobFinish        MessageEvent = "PUBLISHJOBFINISH"           // 发布任务完成
	MessageEventUserInfoModified        MessageEvent = "user_info_modified"         // 用户授权信息变更事件推送

	// 企业微信事件

	MessageWorkTypeChangeContact      MessageEvent = "change_contact"           // 通讯录变更
	MessageEventBatchJobResult        MessageEvent = "batch_job_result"         // 异步任务完成
	MessageEventEnterAgent            MessageEvent = "enter_agent"              // 进入应用
	MessageEventViewWork              MessageEvent = "view"                     // 点击菜单跳转链接时的事件推送
	MessageEventOpenApprovalChange    MessageEvent = "open_approval_change"     // 审批状态通知
	MessageEventShareAgentChange      MessageEvent = "share_agent_change"       // 企业互联共享应用事件回调
	MessageEventShareChainChange      MessageEvent = "share_chain_change"       // 上下游共享应用事件回调
	MessageEventTemplateCardEvent     MessageEvent = "template_card_event"      // 模板卡片事件推送
	MessageEventTemplateCardMenuEvent MessageEvent = "template_card_menu_event" // 通用模板卡片右上角菜单事件推送

	// 开放平台事件

	MessageEventVerifyTicket              MessageEvent = "component_verify_ticket"    // 返回ticket
	MessageEventAuthorized                MessageEvent = "authorized"                 // 授权
	MessageEventUnauthorized              MessageEvent = "unauthorized"               // 取消授权
	MessageEventUpdateAuthorized          MessageEvent = "updateauthorized"           // 更新授权
	MessageEventNotifyThirdFasterRegister MessageEvent = "notify_third_fasteregister" // 注册审核事件推送
)

// CommonMessage
// @Description: 消息中通用的结构
type CommonMessage struct {
	XMLName      xml.Name    `xml:"xml"`
	ToUserName   string      `json:"ToUserName" xml:"ToUserName"`
	FromUserName string      `json:"FromUserName" xml:"FromUserName"`
	CreateTime   int64       `json:"CreateTime" xml:"CreateTime"`
	MsgType      MessageType `json:"MsgType" xml:"MsgType"`
	Nonce        string      `json:"-" xml:"-"`
	ctx          *Context
}

// Message
// @Description: 微信推送的消息
type Message struct {
	CommonMessage
	// 普通消息
	MsgId         int64  `json:"MsgId,omitempty" xml:"MsgId,omitempty"`               // 普通消息的ID
	TemplateMsgID int64  `json:"MsgID,omitempty" xml:"MsgID,omitempty"`               // 模板消息的ID                                             // 模板消息推送成功的消息是MsgID
	Content       string `json:"Content,omitempty" xml:"Content,omitempty"`           // 文本消息
	PicUrl        string `json:"PicUrl,omitempty" xml:"PicUrl,omitempty"`             // 图片消息
	MediaId       string `json:"MediaId,omitempty" xml:"MediaId,omitempty"`           // 图片消息、语音消息、视频消息
	Format        string `json:"Format,omitempty" xml:"Format,omitempty"`             // 语音消息，语音格式，如amr，speex等
	Recognition   string `json:"Recognition,omitempty" xml:"Recognition,omitempty"`   // 语音消息，识别结果
	ThumbMediaId  string `json:"ThumbMediaId,omitempty" xml:"ThumbMediaId,omitempty"` // 视频消息，缩略图
	LocationX     string `json:"Location_X,omitempty" xml:"Location_X,omitempty"`     // 位置消息，纬度
	LocationY     string `json:"Location_Y,omitempty" xml:"Location_Y,omitempty"`     // 位置消息，经度
	Scale         int64  `json:"Scale,omitempty" xml:"Scale,omitempty"`               // 位置消息，地图缩放大小
	Label         string `json:"Label,omitempty" xml:"Label,omitempty"`               // 位置消息，地理位置信息
	Title         string `json:"Title,omitempty" xml:"Title,omitempty"`               // 链接消息，标题
	Description   string `json:"Description,omitempty" xml:"Description,omitempty"`   // 链接消息，描述
	Url           string `json:"Url,omitempty" xml:"Url,omitempty"`                   // 链接消息

	// 事件消息
	Event      MessageEvent `json:"Event,omitempty" xml:"Event,omitempty"`           // 事件消息
	EventKey   string       `json:"EventKey,omitempty" xml:"EventKey,omitempty"`     // 事件，二维码消息、关注、菜单
	Ticket     string       `json:"Ticket,omitempty" xml:"Ticket,omitempty"`         // 事件，二维码消息，二维码ticket
	Latitude   string       `json:"Latitude,omitempty" xml:"Latitude,omitempty"`     // 事件，地理位置，纬度
	Longitude  string       `json:"Longitude,omitempty" xml:"Longitude,omitempty"`   // 事件，地理位置，经度
	Precision  int64        `json:"Precision,omitempty" xml:"Precision,omitempty"`   // 事件，地理位置，精度
	OpenID     string       `json:"OpenID,omitempty" xml:"OpenID,omitempty"`         // 授权用户资料变更,
	AppID      string       `json:"AppID,omitempty" xml:"AppID,omitempty"`           // 授权用户资料变更,公众号的AppID
	RevokeInfo string       `json:"RevokeInfo,omitempty" xml:"RevokeInfo,omitempty"` // 授权用户资料变更,用户撤回的H5授权信息，201:地址,202:发票信息,203:卡券信息,204:麦克风,205:昵称和头像,206:位置信息,207:选中的图片或视频
	TransInfo  *TransInfo   `json:"TransInfo,omitempty" xml:"TransInfo,omitempty"`   // 消息转发到指定客服
	SubscribeMsgPopupEvent
	SubscribeMsgChangeEvent
	SubscribeMsgSentEvent

	// 企业微信
	AgentID          int64             `json:"AgentID,omitempty" xml:"AgentID,omitempty"`                   // 企业应用的ID
	AppType          string            `json:"AppType,omitempty" xml:"AppType,omitempty"`                   // app类型，在企业微信固定返回wxwork，在微信不返回该字段
	BatchJob         *BatchJob         `json:"BatchJob,omitempty" xml:"BatchJob,omitempty"`                 // 异步任务
	ChangeType       string            `json:"ChangeType,omitempty" xml:"ChangeType,omitempty"`             // 通讯录变更类型
	Id               int64             `json:"Id,omitempty" xml:"Id,omitempty"`                             // ID
	Name             string            `json:"Name,omitempty" xml:"Name,omitempty"`                         // 部门名称
	ParentId         int64             `json:"ParentId,omitempty" xml:"ParentId,omitempty"`                 // 父部门ID
	Order            int64             `json:"Order,omitempty" xml:"Order,omitempty"`                       // 排序
	TagId            int64             `json:"TagId,omitempty" xml:"TagId,omitempty"`                       // 标签Id
	AddUserItems     string            `json:"AddUserItems,omitempty" xml:"AddUserItems,omitempty"`         // 标签中新增的成员userid列表，用逗号分隔
	DelUserItems     string            `json:"DelUserItems,omitempty" xml:"DelUserItems,omitempty"`         // 标签中删除的成员userid列表，用逗号分隔
	AddPartyItems    string            `json:"AddPartyItems,omitempty" xml:"AddPartyItems,omitempty"`       // 标签中新增的部门id列表，用逗号分隔
	DelPartyItems    string            `json:"DelPartyItems,omitempty" xml:"DelPartyItems,omitempty"`       // 标签中删除的部门id列表，用逗号分隔
	ScanCodeInfo     *ScanCodeInfo     `json:"ScanCodeInfo,omitempty" xml:"ScanCodeInfo,omitempty"`         // 扫描信息
	SendPicsInfo     *SendPicsInfo     `json:"SendPicsInfo,omitempty" xml:"SendPicsInfo,omitempty"`         // 发送的图片信息
	SendLocationInfo *SendLocationInfo `json:"SendLocationInfo,omitempty" xml:"SendLocationInfo,omitempty"` // 发送的位置信息
	ApprovalInfo     *ApprovalInfo     `json:"ApprovalInfo,omitempty" xml:"ApprovalInfo,omitempty"`         // 审批信息
	TemplateCard
	UserEvent
	CustomerEvent
}
type BatchJob struct {
	JobId   string `json:"JobId,omitempty" xml:"JobId,omitempty"`     // 异步任务ID
	JobType string `json:"JobType,omitempty" xml:"JobType,omitempty"` // 操作类型，字符串，目前分别有：sync_user(增量更新成员)、 replace_user(全量覆盖成员）、invite_user(邀请成员关注）、replace_party(全量覆盖部门)
	ErrCode int64  `json:"ErrCode,omitempty" xml:"ErrCode,omitempty"` // 返回码
	ErrMsg  string `json:"ErrMsg,omitempty" xml:"ErrMsg,omitempty"`   // 返回码描述
}
type ScanCodeInfo struct {
	ScanType   string `json:"ScanType,omitempty" xml:"ScanType,omitempty"`     // 扫描类型，一般是qrcode
	ScanResult string `json:"ScanResult,omitempty" xml:"ScanResult,omitempty"` // 扫描结果，即二维码对应的字符串信息
}
type SendPicsInfo struct {
	Count   int64 `json:"Count,omitempty" xml:"Count,omitempty"` // 发送的图片数量
	PicList []struct {
		PicMd5Sum string `json:"PicMd5Sum,omitempty" xml:"PicMd5Sum,omitempty"` // 图片的MD5值，开发者若需要，可用于验证接收到图片
	} `json:"PicList,omitempty" xml:"PicList>item,omitempty"` // 图片列表
}
type SendLocationInfo struct {
	Location_X string `json:"Location_X,omitempty" xml:"Location_X,omitempty"` // X坐标信息
	Location_Y string `json:"Location_Y,omitempty" xml:"Location_Y,omitempty"` // Y坐标信息
	Scale      string `json:"Scale,omitempty" xml:"Scale,omitempty"`           // 精度，可理解为精度或者比例尺、越精细的话 scale越高
	Label      string `json:"Label,omitempty" xml:"Label,omitempty"`           // 地理位置的字符串信息
	Poiname    string `json:"Poiname,omitempty" xml:"Poiname,omitempty"`       // POI的名字，可能为空
}
type ApprovalInfo struct {
	ThirdNo        string `json:"ThirdNo,omitempty" xml:"ThirdNo,omitempty"`               // 审批单编号，由开发者在发起申请时自定义
	OpenSpName     string `json:"OpenSpName,omitempty" xml:"OpenSpName,omitempty"`         // 审批模板名称
	OpenTemplateId string `json:"OpenTemplateId,omitempty" xml:"OpenTemplateId,omitempty"` // 审批模板id
	OpenSpStatus   int64  `json:"OpenSpStatus,omitempty" xml:"OpenSpStatus,omitempty"`     // 申请单当前审批状态：1-审批中；2-已通过；3-已驳回；4-已取消
	ApplyTime      int64  `json:"ApplyTime,omitempty" xml:"ApplyTime,omitempty"`           // 提交申请时间
	ApplyUserName  string `json:"ApplyUserName,omitempty" xml:"ApplyUserName,omitempty"`   // 提交者姓名
	ApplyUserId    string `json:"ApplyUserId,omitempty" xml:"ApplyUserId,omitempty"`       // 提交者userid
	ApplyUserParty string `json:"ApplyUserParty,omitempty" xml:"ApplyUserParty,omitempty"` // 提交者所在部门
	ApplyUserImage string `json:"ApplyUserImage,omitempty" xml:"ApplyUserImage,omitempty"` // 提交者头像
	Approverstep   int64  `json:"approverstep,omitempty" xml:"approverstep,omitempty"`     // 当前审批节点：0-第一个审批节点；1-第二个审批节点…以此类推
	ApprovalNodes  []struct {
		NodeStatus int64 `json:"NodeStatus,omitempty" xml:"NodeStatus,omitempty"` // 节点审批操作状态：1-审批中；2-已同意；3-已驳回；4-已转审
		NodeAttr   int64 `json:"NodeAttr,omitempty" xml:"NodeAttr,omitempty"`     // 审批节点属性：1-或签；2-会签
		NodeType   int64 `json:"NodeType,omitempty" xml:"NodeType,omitempty"`     // 审批节点类型：1-固定成员；2-标签；3-上级
		Items      []struct {
			ItemName   string `json:"ItemName,omitempty" xml:"ItemName,omitempty"`     // 分支审批人姓名
			ItemUserId string `json:"ItemUserId,omitempty" xml:"ItemUserId,omitempty"` // 分支审批人userid
			ItemImage  string `json:"ItemImage,omitempty" xml:"ItemImage,omitempty"`   // 分支审批人头像
			ItemStatus int64  `json:"ItemStatus,omitempty" xml:"ItemStatus,omitempty"` // 分支审批审批操作状态：1-审批中；2-已同意；3-已驳回；4-已转审
			ItemSpeech string `json:"ItemSpeech,omitempty" xml:"ItemSpeech,omitempty"` // 分支审批人审批意见
			ItemOpTime int64  `json:"ItemOpTime,omitempty" xml:"ItemOpTime,omitempty"` // 分支审批人操作时间
		} `json:"Items,omitempty" xml:"Items>Item,omitempty"` // 审批节点信息，当节点为标签或上级时，一个节点可能有多个分支
	} `json:"ApprovalNodes,omitempty" xml:"ApprovalNodes>ApprovalNode,omitempty"` // 审批流程信息
	NotifyNodes []struct {
		ItemName   string `json:"ItemName,omitempty" xml:"ItemName,omitempty"`     // 抄送人姓名
		ItemUserId string `json:"ItemUserId,omitempty" xml:"ItemUserId,omitempty"` // 抄送人userid
		ItemImage  string `json:"ItemImage,omitempty" xml:"ItemImage,omitempty"`   // 抄送人头像
	} `json:"NotifyNodes,omitempty" xml:"NotifyNodes>NotifyNode,omitempty"` // 抄送信息，可能有多个抄送人
}
type TemplateCard struct {
	TaskId        string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`             // 与发送模板卡片消息时指定的task_id相同
	CardType      string `json:"CardType,omitempty" xml:"CardType,omitempty"`         // 通用模板卡片的类型，类型有"text_notice", "news_notice", "button_interaction", "vote_interaction", "multiple_interaction"五种
	ResponseCode  string `json:"ResponseCode,omitempty" xml:"ResponseCode,omitempty"` // 用于调用更新卡片接口的ResponseCode，24小时内有效，且只能使用一次
	SelectedItems []struct {
		QuestionKey string   `json:"QuestionKey,omitempty" xml:"QuestionKey,omitempty"`      // 问题的key值
		OptionIds   []string `json:"OptionIds,omitempty" xml:"OptionIds>OptionId,omitempty"` // 对应问题的选项列表
	} `json:"SelectedItems,omitempty" xml:"SelectedItems>SelectedItem,omitempty"`
}
type UserEvent struct {
	UserID         string      `json:"UserID,omitempty" xml:"UserID,omitempty"`                 // 成员UserID
	NewUserID      string      `json:"NewUserID,omitempty" xml:"NewUserID,omitempty"`           // 成员新UserID
	Name           string      `json:"Name,omitempty" xml:"Name,omitempty"`                     // 成员名称;代开发自建应用需要管理员授权才返回
	Department     string      `json:"Department,omitempty" xml:"Department,omitempty"`         // 成员部门列表，仅返回该应用有查看权限的部门id
	MainDepartment string      `json:"MainDepartment,omitempty" xml:"MainDepartment,omitempty"` // 主部门
	IsLeaderInDept string      `json:"IsLeaderInDept,omitempty" xml:"IsLeaderInDept,omitempty"` // 表示所在部门是否为部门负责人，0-否，1-是，顺序与Department字段的部门逐一对应。上游共享的应用不返回该字段
	DirectLeader   string      `json:"DirectLeader,omitempty" xml:"DirectLeader,omitempty"`     // 直属上级UserID，最多5个。代开发的自建应用和上游共享的应用不返回该字段
	Position       string      `json:"Position,omitempty" xml:"Position,omitempty"`             // 职位信息。长度为0~64个字节;代开发自建应用需要管理员授权才返回。上游共享的应用不返回该字段
	Mobile         string      `json:"Mobile,omitempty" xml:"Mobile,omitempty"`                 // 手机号码;代开发自建应用需要管理员授权才返回。上游共享的应用不返回该字段
	Gender         int64       `json:"Gender,omitempty" xml:"Gender,omitempty"`                 // 性别，1表示男性，2表示女性。上游共享的应用不返回该字段
	Email          string      `json:"Email,omitempty" xml:"Email,omitempty"`                   // 邮箱;代开发自建应用需要管理员授权才返回。上游共享的应用不返回该字段
	BizMail        string      `json:"BizMail,omitempty" xml:"BizMail,omitempty"`               // 企业邮箱;代开发自建应用不返回该字段。上游共享的应用不返回该字段
	Status         interface{} `json:"Status,omitempty" xml:"Status,omitempty"`                 // 激活状态：1=已激活 2=已禁用 4=未激活 已激活代表已激活企业微信或已关注微信插件（原企业号）5=成员退出
	Avatar         string      `json:"Avatar,omitempty" xml:"Avatar,omitempty"`                 // 头像url。注：如果要获取小图将url最后的”/0”改成”/100”即可。上游共享的应用不返回该字段
	Alias          string      `json:"Alias,omitempty" xml:"Alias,omitempty"`                   // 成员别名。上游共享的应用不返回该字段
	Telephone      string      `json:"Telephone,omitempty" xml:"Telephone,omitempty"`           // 座机;代开发自建应用需要管理员授权才返回。上游共享的应用不返回该字段
	Address        string      `json:"Address,omitempty" xml:"Address,omitempty"`               // 地址;代开发自建应用需要管理员授权才返回。上游共享的应用不返回该字段
	ExtAttr        []struct {
		Name string   `json:"Name,omitempty" xml:"Name,omitempty"`
		Type int64    `json:"Type,omitempty" xml:"Type,omitempty"`
		Text []string `json:"Text,omitempty" xml:"Text>Value,omitempty"`
		Web  *Web     `json:"Web,omitempty" xml:"Web,omitempty"`
	} `json:"ExtAttr,omitempty" xml:"ExtAttr>Item,omitempty"` // 扩展属性;代开发自建应用需要管理员授权才返回。上游共享的应用不返回该字段
}
type Web struct {
	Title string `json:"Title,omitempty" xml:"Title,omitempty"`
	Url   string `json:"Url,omitempty" xml:"Url,omitempty"`
}
type CustomerEvent struct {
	ExternalUserID string `json:"ExternalUserID,omitempty" xml:"ExternalUserID,omitempty"` // 外部联系人的userid，注意不是企业成员的帐号
	State          string `json:"State,omitempty" xml:"State,omitempty"`                   // 添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道
	WelcomeCode    string `json:"WelcomeCode,omitempty" xml:"WelcomeCode,omitempty"`       // 欢迎语code，可用于发送欢迎语
	Source         string `json:"Source,omitempty" xml:"Source,omitempty"`                 // 删除客户的操作来源，DELETE_BY_TRANSFER表示此客户是因在职继承自动被转接成员删除
	FailReason     string `json:"FailReason,omitempty" xml:"FailReason,omitempty"`         // 接替失败的原因, customer_refused-客户拒绝， customer_limit_exceed-接替成员的客户数达到上限
	ChatId         string `json:"ChatId,omitempty" xml:"ChatId,omitempty"`                 // 群ID
	UpdateDetail   string `json:"UpdateDetail,omitempty" xml:"UpdateDetail,omitempty"`     // 变更详情
	JoinScene      int64  `json:"JoinScene,omitempty" xml:"JoinScene,omitempty"`           // 当是成员入群时有值。表示成员的入群方式
	QuitScene      int64  `json:"QuitScene,omitempty" xml:"QuitScene,omitempty"`           // 当是成员退群时有值。表示成员的退群方式
	MemChangeCnt   int64  `json:"MemChangeCnt,omitempty" xml:"MemChangeCnt,omitempty"`     // 当是成员入群或退群时有值。表示成员变更数量
	StrategyId     string `json:"StrategyId,omitempty" xml:"StrategyId,omitempty"`         // 标签或标签组所属的规则组id，只回调给“客户联系”应用
}
type SubscribeMsgPopupEvent struct {
	SubscribeMsgPopupEvent []struct {
		TemplateId            string `json:"TemplateId,omitempty" xml:"TemplateId,omitempty"`                       // 模板 id（一次订阅可能有多条通知，带有多个 id）
		SubscribeStatusString string `json:"SubscribeStatusString,omitempty" xml:"SubscribeStatusString,omitempty"` // 用户点击行为（同意、取消发送通知）
		PopupScene            int64  `json:"PopupScene,omitempty" xml:"PopupScene,omitempty"`                       // 1弹窗来自H5页面,2弹窗来自图文消息
	} `json:"SubscribeMsgPopupEvent,omitempty" xml:"SubscribeMsgPopupEvent>List,omitempty"` // 用户操作订阅通知弹窗
}

type SubscribeMsgChangeEvent struct {
	SubscribeMsgChangeEvent []struct {
		TemplateId            string `json:"TemplateId,omitempty" xml:"TemplateId,omitempty"`                       // 模板 id（一次订阅可能有多条通知，带有多个 id）
		SubscribeStatusString string `json:"SubscribeStatusString,omitempty" xml:"SubscribeStatusString,omitempty"` // 用户点击行为（仅推送用户拒收通知）
	} `json:"SubscribeMsgChangeEvent,omitempty" xml:"SubscribeMsgChangeEvent>List,omitempty"` // 用户管理订阅通知
}
type SubscribeMsgSentEvent struct {
	SubscribeMsgSentEvent []struct {
		TemplateId  string `json:"TemplateId,omitempty" xml:"TemplateId,omitempty"`   // 模板 id（一次订阅可能有多条通知，带有多个 id）
		MsgID       string `json:"MsgID,omitempty" xml:"MsgID,omitempty"`             // 消息 id
		ErrorCode   int64  `json:"ErrorCode,omitempty" xml:"ErrorCode,omitempty"`     // 推送结果状态码（0表示成功）
		ErrorStatus string `json:"ErrorStatus,omitempty" xml:"ErrorStatus,omitempty"` // 推送结果状态码文字含义
	} `json:"SubscribeMsgSentEvent,omitempty" xml:"SubscribeMsgSentEvent>List,omitempty"` // 发送订阅通知
}
type TransInfo struct {
	KfAccount string `json:"KfAccount,omitempty" xml:"KfAccount,omitempty"` // 指定会话接入的客服账号
}

// MessageReply
// @Description: 回复消息
type MessageReply struct {
	CommonMessage
	// 回复文本
	Content string `xml:"Content,omitempty"`
	// 回复图片
	Image *MediaID `xml:"Image,omitempty"`
	// 回复录音
	Voice *MediaID `xml:"Voice,omitempty"`
	// 回复视频
	Video *Video `xml:"Video,omitempty"`
	// 回复音乐
	Music *Music `xml:"Music,omitempty"`
	// 回复图文
	ArticleCount int `xml:"ArticleCount,omitempty"`
}
type MessageReplyArticles struct {
	MessageReply
	Articles []*Article `xml:"Articles>item,omitempty"`
}

type MediaID struct {
	MediaID string `xml:"MediaId,omitempty"`
}
type Video struct {
	MediaID     string `xml:"MediaId,omitempty"`
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
}
type Music struct {
	Title        string `xml:"Title,omitempty"`
	Description  string `xml:"Description,omitempty"`
	MusicURL     string `xml:"MusicUrl,omitempty"`
	HQMusicURL   string `xml:"HQMusicUrl,omitempty"`
	ThumbMediaID string `xml:"ThumbMediaId,omitempty"`
}
type Article struct {
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
	PicURL      string `xml:"PicUrl,omitempty"`
	URL         string `xml:"Url,omitempty"`
}
