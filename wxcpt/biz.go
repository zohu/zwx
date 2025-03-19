package wxcpt

import "encoding/xml"

type BizMsgSendJson struct {
	Encrypt      string `json:"encrypt"`
	Msgsignature string `json:"msgsignature"`
	Nonce        string `json:"nonce"`
	Timestamp    string `json:"timestamp"`
}

func NewBizMsgSendJson(encrypt, signature, timestamp, nonce string) *BizMsgSendJson {
	return &BizMsgSendJson{
		Encrypt:      encrypt,
		Msgsignature: signature,
		Nonce:        nonce,
		Timestamp:    timestamp,
	}
}

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

type BizMsgSendXml struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATA    `xml:"Encrypt"`
	Msgsignature CDATA    `xml:"MsgSignature"`
	Nonce        CDATA    `xml:"Nonce"`
	Timestamp    string   `xml:"TimeStamp"`
}

func NewBizMsgSendXml(encrypt, signature, timestamp, nonce string) *BizMsgSendXml {
	return &BizMsgSendXml{
		Encrypt:      CDATA(encrypt),
		Msgsignature: CDATA(signature),
		Nonce:        CDATA(nonce),
		Timestamp:    timestamp,
	}
}
