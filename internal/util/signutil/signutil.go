package signutil

import (
	"strings"
	"github.com/fenfenbingo/bingosession/internal/util/convert"
	"sort"
	"github.com/fenfenbingo/bingosession/internal/util"
)

// A data structure to hold a key/value pair.
type pair struct {
	Key   string
	Value string
}

type pairList []pair

func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return strings.Compare(p[i].Key, p[j].Key) < 0 }

//getSign generate the signature
//getSign用于生成签名
func GetSign(m map[string]interface{}, signKey string) string {
	p := make(pairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = pair{k, convert.AnyToStr(v)}
		i++
	}
	sort.Sort(p)
	signStr := ""
	for _, v := range p {
		signStr += ("&" + v.Key + "=" + v.Value)
	}
	return util.MD5([]byte(signStr + signKey))
}

//checkSign check if the signature is valid
//checkSign用于检查签名是否有效
func CheckSign(m map[string]interface{}, sign string, signKey string) bool {
	return GetSign(m, signKey) == sign
}
