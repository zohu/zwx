package zwx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	"github.com/zohu/zwx/utils"
	"strings"
)

type Api string

func (a Api) WithPath(path string) string {
	return fmt.Sprintf(
		"%s/%s",
		strings.TrimSuffix(a.String(), "/"),
		strings.TrimPrefix(path, "/"),
	)
}
func (a Api) String() string {
	return string(a)
}

const (
	ApiCgiBin     Api = "https://api.weixin.qq.com/cgi-bin"
	ApiMpCgiBin   Api = "https://mp.weixin.qq.com/cgi-bin"
	ApiWorkCgiBin Api = "https://qyapi.weixin.qq.com/cgi-bin"
	ApiWxa        Api = "https://api.weixin.qq.com/wxa"
	ApiWxaapi     Api = "https://api.weixin.qq.com/wxaapi"
	ApiSns        Api = "https://api.weixin.qq.com/sns"
)

type Method string

func (m Method) String() string {
	return string(m)
}

const (
	MethodGet  Method = fasthttp.MethodGet
	MethodPost Method = fasthttp.MethodPost
)

type Http struct {
	c        *fasthttp.Client
	req      *fasthttp.Request
	resp     *fasthttp.Response
	handlers []func(resp *fasthttp.Response)
	errs     []string
	debug    bool
	logger   Logger
}

func NewHttp(method Method, uri string) *Http {
	h := &Http{
		c:    &fasthttp.Client{},
		req:  fasthttp.AcquireRequest(),
		resp: fasthttp.AcquireResponse(),
	}
	h.req.SetRequestURI(uri)
	h.req.Header.SetMethod(method.String())
	return h
}
func (h *Http) SetJson(body any) *Http {
	d, err := sonic.Marshal(body)
	if err != nil {
		h.errs = append(h.errs, fmt.Sprintf("body marshal json error: %v", err))
		return h
	}
	h.req.Header.SetContentType("application/json")
	h.req.SetBody(d)
	return h
}
func (h *Http) SetXml(body any) *Http {
	d, err := xml.Marshal(body)
	if err != nil {
		h.errs = append(h.errs, fmt.Sprintf("body marshal xml error: %v", err))
		return h
	}
	h.req.Header.SetContentType("application/xml")
	h.req.SetBody(d)
	return h
}
func (h *Http) SetForm(body any) *Http {
	d, err := sonic.Marshal(body)
	if err != nil {
		h.errs = append(h.errs, fmt.Sprintf("body marshal form error: %v", err))
		return h
	}
	h.req.Header.SetContentType("application/x-www-form-urlencoded")
	h.req.SetBody(d)
	return h
}
func (h *Http) SetAccessToken(token string) *Http {
	h.req.URI().QueryArgs().Add("access_token", token)
	return h
}
func (h *Http) SetQuery(query map[string]string) *Http {
	for k, v := range query {
		h.req.URI().QueryArgs().Add(k, fmt.Sprintf("%v", v))
	}
	return h
}
func (h *Http) SetHeader(header map[string]string) *Http {
	for k, v := range header {
		h.req.Header.Set(k, v)
	}
	return h
}
func (h *Http) BindJson(v any) *Http {
	h.handlers = append(h.handlers, func(resp *fasthttp.Response) {
		if err := sonic.Unmarshal(resp.Body(), v); err != nil {
			h.errs = append(h.errs, fmt.Sprintf("resp unmarshal json error: %v", err))
		}
	})
	return h
}
func (h *Http) BindXml(v any) *Http {
	h.handlers = append(h.handlers, func(resp *fasthttp.Response) {
		if err := xml.Unmarshal(resp.Body(), v); err != nil {
			h.errs = append(h.errs, fmt.Sprintf("resp unmarshal xml error: %v", err))
		}
	})
	return h
}
func (h *Http) BindJsonOrBytes(obj any, v *[]byte) *Http {
	h.handlers = append(h.handlers, func(resp *fasthttp.Response) {
		if err := sonic.Unmarshal(resp.Body(), obj); err == nil {
			return
		}
		*v = resp.Body()
	})
	return h
}
func (h *Http) Debug(debug bool, logger Logger) *Http {
	h.debug = debug
	h.logger = logger
	return h
}
func (h *Http) Do() error {
	defer fasthttp.ReleaseRequest(h.req)
	defer fasthttp.ReleaseResponse(h.resp)
	// 发送请求
	if err := h.c.Do(h.req, h.resp); err != nil {
		h.errs = append(h.errs, fmt.Sprintf("request error: %v", err))
	} else {
		// 序列化返回值
		for _, handler := range h.handlers {
			handler(h.resp)
		}
	}
	// 是否debug
	if h.debug {
		buf := utils.NewBuffer()
		_ = buf.WriteByte('\n')
		_, _ = buf.Write(h.req.Header.Method())
		_ = buf.WriteByte(' ')
		_, _ = buf.Write(h.req.URI().Scheme())
		_, _ = buf.WriteString("//")
		_, _ = buf.Write(h.req.URI().Host())
		_, _ = buf.Write(h.req.URI().Path())
		_ = buf.WriteByte(' ')
		_, _ = buf.Write(h.req.Header.Protocol())
		_ = buf.WriteByte('\n')
		_, _ = buf.WriteString("Header:\n")
		h.req.Header.VisitAll(func(k, v []byte) {
			_ = buf.WriteByte('\t')
			_, _ = buf.Write(k)
			_, _ = buf.WriteString(": ")
			_, _ = buf.Write(v)
			_ = buf.WriteByte('\n')
		})
		qs := h.req.URI().QueryArgs()
		if qs.Len() > 0 {
			_, _ = buf.WriteString("Query:\n")
			qs.VisitAll(func(k, v []byte) {
				_ = buf.WriteByte('\t')
				_, _ = buf.Write(k)
				_, _ = buf.WriteString(": ")
				_, _ = buf.Write(v)
				_ = buf.WriteByte('\n')
			})
		}
		bd := h.req.Body()
		if len(bd) > 0 {
			_, _ = buf.WriteString("Body:\n\t")
			_, _ = buf.Write(bd)
			_ = buf.WriteByte('\n')
		}
		_, _ = buf.WriteString("Response:\n\t")
		_, _ = buf.Write(h.resp.Body())
		_ = buf.WriteByte('\n')
		h.logger.Debugf(buf.String())
	}
	// 检查是否有错误
	if len(h.errs) > 0 {
		return errors.New(strings.Join(h.errs, "\n"))
	}
	return nil
}
