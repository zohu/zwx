package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

// RandomIntStr
// @Description: 生成随机数字符串
// @receiver o
// @return string
func RandomIntStr(len int) string {
	return fmt.Sprintf("%0*d", len, rand.Intn(int(math.Pow10(len))))
}

// RandomStr
// @Description: 随机生成字符串
// @param length
// @return string
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Signature
// @Description: sha1签名
// @param params
// @return string
func Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		_, _ = io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Md5
// @Auth: oak  2021-10-15 18:43:27
// @Description:  MD5
// @param v
// @return string
func Md5(v string) string {
	data := []byte(v)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func StructToMap(obj interface{}) map[string]string {
	objMap := make(map[string]string)
	d, _ := sonic.Marshal(obj)
	_ = sonic.Unmarshal(d, &objMap)
	return objMap
}

func FirstTruth[T any](args ...T) T {
	for _, item := range args {
		if !reflect.ValueOf(item).IsZero() {
			return item
		}
	}
	return args[0]
}
