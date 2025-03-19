package zwx

type Prefix string

const (
	PrefixAppList Prefix = "wx:1"
	PrefixApp     Prefix = "wx:2"
	PrefixRetry   Prefix = "wx:3"
)

func (p Prefix) Key(val ...string) string {
	key := string(p)
	for _, k := range val {
		key = key + ":" + k
	}
	return key
}

type WxResponse struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}
