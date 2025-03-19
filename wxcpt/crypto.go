package wxcpt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/zohu/zwx/utils"
	"sort"
	"strings"
)

type BizMsgRecv struct {
	Tousername string `json:"tousername" xml:"ToUserName"`
	Encrypt    string `json:"encrypt" xml:"Encrypt"`
	Agentid    string `json:"agentid,omitempty" xml:"AgentID,omitempty"`
}
type BizMsgCrypt struct {
	Token          string
	EncodingAesKey string
	Appid          string
}

func NewBizMsgCrypt(token, encodingAeskey, appid string) *BizMsgCrypt {
	return &BizMsgCrypt{
		Token:          token,
		EncodingAesKey: encodingAeskey + "=",
		Appid:          appid,
	}
}

// EncryptJsonMsg
// @Description: json消息加密
// @receiver cpt
// @param replyMsg
// @param timestamp
// @param nonce
// @return *JsonBizMsg4Send
// @return error
func (cpt *BizMsgCrypt) EncryptJsonMsg(replyMsg, timestamp, nonce string) (*BizMsgSendJson, error) {
	randStr := utils.RandomStr(16)
	var buffer bytes.Buffer
	buffer.WriteString(randStr)

	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(replyMsg)))
	buffer.Write(msgLenBuf)
	buffer.WriteString(replyMsg)
	buffer.WriteString(cpt.Appid)

	tmpCiphertext, err := cpt.cbcEncrypt(buffer.String())
	if nil != err {
		return nil, err
	}
	ciphertext := string(tmpCiphertext)
	signature := cpt.calSignature(timestamp, nonce, ciphertext)
	return NewBizMsgSendJson(ciphertext, signature, timestamp, nonce), nil
}

// EncryptXmlMsg
// @Description: xml消息加密
// @receiver cpt
// @param replyMsg
// @param timestamp
// @param nonce
// @return *XmlBizMsg4Send
// @return error
func (cpt *BizMsgCrypt) EncryptXmlMsg(replyMsg, timestamp, nonce string) (*BizMsgSendXml, error) {
	randStr := utils.RandomStr(16)
	var buffer bytes.Buffer
	buffer.WriteString(randStr)

	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(replyMsg)))
	buffer.Write(msgLenBuf)
	buffer.WriteString(replyMsg)
	buffer.WriteString(cpt.Appid)

	tmpCiphertext, err := cpt.cbcEncrypt(buffer.String())
	if nil != err {
		return nil, err
	}
	ciphertext := string(tmpCiphertext)
	signature := cpt.calSignature(timestamp, nonce, ciphertext)
	return NewBizMsgSendXml(ciphertext, signature, timestamp, nonce), nil
}

// DecryptMsg
// @Description: 消息解密
// @receiver cpt
// @param msgSignature
// @param timestamp
// @param nonce
// @param msg4Recv
// @param protocolType
// @return []byte
// @return *CryptError
func (cpt *BizMsgCrypt) DecryptMsg(msgSignature, timestamp, nonce string, msg4Recv *BizMsgRecv) ([]byte, error) {
	signature := cpt.calSignature(timestamp, nonce, msg4Recv.Encrypt)
	if strings.Compare(signature, msgSignature) != 0 {
		return nil, errors.New("signature not equal")
	}
	plaintext, err := cpt.cbcDecrypt(msg4Recv.Encrypt)
	if nil != err {
		return nil, err
	}
	_, _, msg, appid, err := cpt.ParsePlainText(plaintext)
	if nil != err {
		return nil, err
	}
	if len(cpt.Appid) > 0 && strings.Compare(string(appid), cpt.Appid) != 0 {
		return nil, errors.New("appid is not equal")
	}
	return msg, nil
}

// DecryptMsgFromBinary
// @Description: 从二进制解密消息
// @receiver cpt
// @param msgSignature
// @param timestamp
// @param nonce
// @param msg
// @return []byte
// @return error
func (cpt *BizMsgCrypt) DecryptMsgFromBinary(msgSignature, timestamp, nonce string, msg []byte) ([]byte, error) {
	recv := new(BizMsgRecv)
	if err := xml.Unmarshal(msg, recv); err != nil {
		return nil, err
	}
	return cpt.DecryptMsg(msgSignature, timestamp, nonce, recv)
}

func (cpt *BizMsgCrypt) pKCS7Padding(plaintext string, blockSize int) []byte {
	padding := blockSize - (len(plaintext) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	var buffer bytes.Buffer
	buffer.WriteString(plaintext)
	buffer.Write(padText)
	return buffer.Bytes()
}

func (cpt *BizMsgCrypt) pKCS7UnPadding(plaintext []byte, blockSize int) ([]byte, error) {
	plaintextLen := len(plaintext)
	if nil == plaintext || plaintextLen == 0 {
		return nil, errors.New("pKCS7UnPadding error nil or zero")
	}
	if plaintextLen%blockSize != 0 {
		return nil, errors.New("pKCS7UnPadding text not a multiple of the block size")
	}
	paddingLen := int(plaintext[plaintextLen-1])
	return plaintext[:plaintextLen-paddingLen], nil
}

func (cpt *BizMsgCrypt) cbcEncrypt(plaintext string) ([]byte, error) {
	aeskey, err := base64.StdEncoding.DecodeString(cpt.EncodingAesKey)
	if nil != err {
		return nil, errors.New("base64 error")
	}
	const blockSize = 32
	padMsg := cpt.pKCS7Padding(plaintext, blockSize)

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, errors.New("aes error")
	}

	ciphertext := make([]byte, len(padMsg))
	iv := aeskey[:aes.BlockSize]

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, padMsg)
	base64Msg := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(base64Msg, ciphertext)

	return base64Msg, nil
}

func (cpt *BizMsgCrypt) cbcDecrypt(base64EncryptMsg string) ([]byte, error) {
	aeskey, err := base64.StdEncoding.DecodeString(cpt.EncodingAesKey)
	if nil != err {
		return nil, errors.New("base64 decode error")
	}
	encryptMsg, err := base64.StdEncoding.DecodeString(base64EncryptMsg)
	if nil != err {
		return nil, errors.New("base64 decode error")
	}
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, errors.New("aes error")
	}
	if len(encryptMsg) < aes.BlockSize {
		return nil, errors.New("encryptMsg size is not valid")
	}
	iv := aeskey[:aes.BlockSize]
	if len(encryptMsg)%aes.BlockSize != 0 {
		return nil, errors.New("encryptMsg not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptMsg, encryptMsg)
	return encryptMsg, nil
}

func (cpt *BizMsgCrypt) calSignature(timestamp, nonce, data string) string {
	sortArr := []string{cpt.Token, timestamp, nonce, data}
	sort.Strings(sortArr)
	var buffer bytes.Buffer
	for _, value := range sortArr {
		buffer.WriteString(value)
	}
	sha := sha1.New()
	sha.Write(buffer.Bytes())
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func (cpt *BizMsgCrypt) ParsePlainText(plaintext []byte) ([]byte, uint32, []byte, []byte, error) {
	const blockSize = 32
	plaintext, err := cpt.pKCS7UnPadding(plaintext, blockSize)
	if nil != err {
		return nil, 0, nil, nil, err
	}
	textLen := uint32(len(plaintext))
	if textLen < 20 {
		return nil, 0, nil, nil, errors.New("plain is to small 1")
	}
	random := plaintext[:16]
	msgLen := binary.BigEndian.Uint32(plaintext[16:20])
	if textLen < (20 + msgLen) {
		return nil, 0, nil, nil, errors.New("plain is to small 2")
	}
	msg := plaintext[20 : 20+msgLen]
	appid := plaintext[20+msgLen:]
	return random, msgLen, msg, appid, nil
}

func (cpt *BizMsgCrypt) VerifyURL(msgSignature, timestamp, nonce, echostr string) ([]byte, error) {
	signature := cpt.calSignature(timestamp, nonce, echostr)
	if strings.Compare(signature, msgSignature) != 0 {
		return nil, errors.New("signature not equal")
	}
	plaintext, err := cpt.cbcDecrypt(echostr)
	if nil != err {
		return nil, err
	}
	_, _, msg, appid, err := cpt.ParsePlainText(plaintext)
	if nil != err {
		return nil, err
	}
	if len(cpt.Appid) > 0 && strings.Compare(string(appid), cpt.Appid) != 0 {
		return nil, errors.New("appid is not equal")
	}
	return msg, nil
}

// HmacSha256
// @Description: 计算HmacSha256
// @param key
// @param data
// @return []byte
func HmacSha256(key string, data string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))
	return mac.Sum(nil)
}

// HmacSha256ToHex
// @Description: 将加密后的二进制转16进制字符串
// @param key
// @param data
// @return string
func HmacSha256ToHex(key string, data string) string {
	return hex.EncodeToString(HmacSha256(key, data))
}

// HmacSha256ToBase64
// @Description: 将加密后的二进制转Base64字符串
// @param key
// @param data
// @return string
func HmacSha256ToBase64(key string, data string) string {
	return base64.URLEncoding.EncodeToString(HmacSha256(key, data))
}
