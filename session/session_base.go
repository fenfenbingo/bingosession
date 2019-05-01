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

	Config *ProviderConf
}

func (s *sessionBase) GetSessionId() string {
	return s.ID
}

func (s *sessionBase) SetSessionId(sid string) {
	s.ID = sid
}

func (s *sessionBase) Get(key string) interface{} {
	if value, ok := s.Values[key]; ok {
		return value
	}
	return nil
}

func (s *sessionBase) Set(key string, value interface{}) {
	s.Values[key] = value
}

func (s *sessionBase) Delete(key string) {
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



