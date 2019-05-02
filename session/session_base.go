package session

import (
	"net/http"
	"github.com/json-iterator/go"
)

//##############provider################

type providerBase struct {
	StoreUrl string
	Config   *ProviderConf
}

func (p *providerBase) GetCookieContent(r *http.Request) (sid string) {
	c, err := r.Cookie(p.Config.SessCookieName)
	if err != nil || c.Value == "" {
		return ""
	}
	return c.Value
}

func (p *providerBase) providerInit(conf *ProviderConf) error {
	p.StoreUrl = conf.StoreUrl
	p.Config = conf
	return nil
}

//################store######################

type sessionBase struct {
	Req *http.Request

	RespWriter http.ResponseWriter

	Values map[string]interface{}

	ID string
	//不建议可能存在异步操作session的情况，如有可能，请自行实现并发安全
	//Lock sync.RWMutex

	Config *ProviderConf
}

func (s *sessionBase) GetSessionId() string {
	return s.ID
}

func (s *sessionBase) SetSessionId(sid string) {
	s.ID = sid
}

func (s *sessionBase) Get(key string) interface{} {
	//self.Lock.RLock()
	//defer self.Lock.RUnlock()
	if value, ok := s.Values[key]; ok {
		return value
	}
	return nil
}

func (s *sessionBase) Set(key string, value interface{}) {
	//self.Lock.Lock()
	//defer self.Lock.Unlock()
	s.Values[key] = value
}

func (s *sessionBase) Delete(key string) {
	//self.Lock.Lock()
	//defer self.Lock.Unlock()
	delete(s.Values, key)
}

func  (s *sessionBase) Destroy() (err error) {
	s.ID=""
	s.SetCookie("")
	return
}


func (s *sessionBase) LoadObject(ref interface{}) {
	b, _ := jsoniter.Marshal(s.Values)
	jsoniter.Unmarshal(b, ref)
}

func (s *sessionBase) SetCookie(content string) {
	config := s.Config
	maxAge := config.CookieMaxAge
	if content == "" {
		maxAge = -1
	}
	http.SetCookie(s.RespWriter, &http.Cookie{
		Name:     config.SessCookieName,
		Value:    content,
		MaxAge:   maxAge,
		Path:     config.Path,
		Domain:   config.Domain,
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
	})
}



