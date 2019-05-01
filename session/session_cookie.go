package session

import (
	"net/http"
	"github.com/json-iterator/go"
	"encoding/base64"
	"time"
	"github.com/fenfenbingo/bingosession/internal/util"
	"github.com/fenfenbingo/bingosession/internal/util/signutil"
)

//session_cookie 使用cookie存放数据,有效期存放在cookie中（无法即时回收令牌）
//########### provider ##############

type providerCookie struct {
	providerBase
}

func (p *providerCookie) SessionStart(req *http.Request, w http.ResponseWriter) (res ISession, ecode *ErrInfo) {
	return p.sessionStart(req, w)
}

func (p *providerCookie) sessionStart(req *http.Request, w http.ResponseWriter) (res *sessionCookie, ecode *ErrInfo) {
	content := p.GetCookieContent(req)
	sess := &sessionCookie{
		sessionBase{
			Values:     make(map[string]interface{}),
			Req:        req,
			RespWriter: w,
			ID:         "",
			Config:     p.Config,
		},
	}
	if content == "" {
		return sess, ErrEmpty
	}
	//destroy the session if it is invalid
	defer func() {
		if ecode != NoErr {
			res.Destroy()
		}
	}()
	bytes, err := base64.URLEncoding.DecodeString(content)
	if err != nil {
		return sess, ErrInValidSessId
	}
	err = jsoniter.Unmarshal(bytes, &sess.Values)
	if err != nil {
		return sess, ErrCacheFormat
	}
	if v, ok := sess.Values["_uuid"]; ok {
		sess.ID = jsoniter.Wrap(v).ToString()
	}
	if sess.ID == "" {
		return sess, ErrCacheFormat
	}
	//check signature
	sign := ""
	if v, ok := sess.Values["_sign"]; ok {
		switch v.(type) {
		case string:
			{
				sign = v.(string)
				break
			}
		default:
			return sess, ErrCacheFormat
		}
		delete(sess.Values, "_sign")

		if !signutil.CheckSign(sess.Values, sign, p.Config.SignKey) {
			return sess, ErrInvalidSign
		}
	} else {
		return sess, ErrCacheFormat
	}
	var expire int64 = 0
	//check expire
	if v, ok := sess.Values["_expire"]; ok {
		switch v.(type) {
		case float64:
			{
				expire = int64(v.(float64))
				break
			}
		default:
			return sess, ErrCacheFormat
		}
		if time.Now().Unix() >= expire {
			return sess, ErrExpire
		}
	} else {
		return sess, ErrCacheFormat
	}
	return sess, NoErr
}

func (p *providerCookie) SessionDestroy(sid string) {
	//此类型的session的内容是存在cookie里的，类似jwt不适合用于会话管理。看session_cookie_redis.go如何处理
	panic("Don't call this function,use 'ISession.Destroy' instead")
}

//########### store ##############

type sessionCookie struct {
	sessionBase
}

func (s *sessionCookie) Save() (err error) {
	if len(s.Values) == 0 {
		s.Destroy()
		return nil
	}
	if s.ID == "" {
		s.ID = util.GenerateUUID()
	}
	s.Values["_uuid"] = s.ID
	expire := time.Now().Unix() + int64(s.Config.MaxLifeTime)
	s.Values["_expire"] = expire
	s.Values["_sign"] = signutil.GetSign(s.Values, s.Config.SignKey)
	data, err := jsoniter.Marshal(s.Values)
	if err != nil {
		return err
	}
	s.SetCookie(base64.URLEncoding.EncodeToString(data))
	return nil
}

func (s *sessionCookie) SaveObject(v interface{}) (err error) {
	b, err := jsoniter.Marshal(v)
	if err != nil {
		return err
	}
	jsoniter.Unmarshal(b, &s.Values)
	return s.Save()
}

func (s *sessionCookie) Destroy() error {
	s.sessionBase.Destroy()
	return nil
}
