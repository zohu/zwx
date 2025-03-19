package wxnotify

import (
	"encoding/xml"
	"github.com/zohu/zwx/wxcpt"
	"strconv"
	"time"
)

func encryptMsg(ctx *Context, data []byte, timestamp int64, nonce string) *wxcpt.BizMsgSendXml {
	cpt := wxcpt.NewBizMsgCrypt(ctx.NotifyToken(), ctx.NotifyEncodingAesKey(), ctx.AppidMain())
	send, err := cpt.EncryptXmlMsg(string(data), strconv.FormatInt(timestamp, 10), nonce)
	if err != nil {
		ctx.Logger().Errorf("加密失败 %s", err.Error())
		return nil
	}
	return send
}

func (msg *MessageReply) Encrypted() *wxcpt.BizMsgSendXml {
	str, _ := xml.Marshal(msg)
	return encryptMsg(msg.ctx, str, msg.CreateTime, msg.Nonce)
}
func (msg *MessageReplyArticles) Encrypted() *wxcpt.BizMsgSendXml {
	str, _ := xml.Marshal(msg)
	return encryptMsg(msg.ctx, str, msg.CreateTime, msg.Nonce)
}

// ReplyText
// @Description: 回复文本消息
// @receiver ctx
// @param content
// @return *MessageText
func (msg *Message) ReplyText(content string) *MessageReply {
	text := new(MessageReply)
	text.Nonce = msg.Nonce
	text.ctx = msg.ctx
	text.MsgType = MessageTypeText
	text.FromUserName = msg.ToUserName
	text.ToUserName = msg.FromUserName
	text.CreateTime = time.Now().Unix()
	text.Content = content
	return text
}

// ReplyImage
// @Description: 回复图片消息
// @receiver ctx
// @param mediaID
// @return *MessageImage
func (msg *Message) ReplyImage(mediaID string) *MessageReply {
	image := new(MessageReply)
	image.Nonce = msg.Nonce
	image.ctx = msg.ctx
	image.MsgType = MessageTypeImage
	image.FromUserName = msg.ToUserName
	image.ToUserName = msg.FromUserName
	image.CreateTime = time.Now().Unix()
	image.Image.MediaID = mediaID
	return image
}

// ReplyVoice
// @Description: 回复语音消息
// @receiver ctx
// @param mediaID
// @return *MessageVoice
func (msg *Message) ReplyVoice(mediaID string) *MessageReply {
	voice := new(MessageReply)
	voice.Nonce = msg.Nonce
	voice.ctx = msg.ctx
	voice.MsgType = MessageTypeVoice
	voice.FromUserName = msg.ToUserName
	voice.ToUserName = msg.FromUserName
	voice.CreateTime = time.Now().Unix()
	voice.Voice.MediaID = mediaID
	return voice
}

// ReplyVideo
// @Description: 回复视频消息
// @receiver ctx
// @param mediaID
// @param title
// @param description
// @return *MessageVideo
func (msg *Message) ReplyVideo(mediaID, title, description string) *MessageReply {
	video := new(MessageReply)
	video.Nonce = msg.Nonce
	video.ctx = msg.ctx
	video.MsgType = MessageTypeVideo
	video.FromUserName = msg.ToUserName
	video.ToUserName = msg.FromUserName
	video.CreateTime = time.Now().Unix()
	video.Video.MediaID = mediaID
	video.Video.Title = title
	video.Video.Description = description
	return video
}

// ReplyMusic
// @Description: 回复音乐消息
// @receiver ctx
// @param title
// @param description
// @param musicURL
// @param hQMusicURL
// @param thumbMediaID
// @return *MessageMusic
func (msg *Message) ReplyMusic(title, description, musicURL, hQMusicURL, thumbMediaID string) *MessageReply {
	music := new(MessageReply)
	music.Nonce = msg.Nonce
	music.ctx = msg.ctx
	music.MsgType = MessageTypeMusic
	music.FromUserName = msg.ToUserName
	music.ToUserName = msg.FromUserName
	music.CreateTime = time.Now().Unix()
	music.Music.Title = title
	music.Music.Description = description
	music.Music.MusicURL = musicURL
	music.Music.HQMusicURL = hQMusicURL
	music.Music.ThumbMediaID = thumbMediaID
	return music
}

// ReplyNews
// @Description: 回复图文消息
// @receiver ctx
// @param articles
// @return *MessageNews
func (msg *Message) ReplyNews(articles []*Article) *MessageReplyArticles {
	news := new(MessageReplyArticles)
	news.Nonce = msg.Nonce
	news.ctx = msg.ctx
	news.MsgType = MessageTypeNews
	news.FromUserName = msg.ToUserName
	news.ToUserName = msg.FromUserName
	news.CreateTime = time.Now().Unix()
	news.ArticleCount = len(articles)
	news.Articles = articles
	return news
}
