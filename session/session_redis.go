package session

import (
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"net/http"
)

//session_redis 使用redis存放数据

//########## provider #############

type providerRedis struct {
	providerBase
	cPool *redis.Pool
}

func (p *providerRedis) SessionStart(req *http.Request, w http.ResponseWriter) (res ISession, ecode *ErrInfo) {
	sid := p.GetCookieContent(req)
	sess := &sessionRedis{
		sessionBase{
			Values:     make(map[string]interface{}),
			Req:        req,
			RespWriter: w,
			ID:         sid,
			Config:     p.Config,
		},
		p.cPool,
	}
	if sid == "" {
		return sess, ErrEmpty
	}
	defer func() {
		if ecode != NoErr {
			res.SetSessionId("")
		}
	}()
	if len(sid) != 32 {
		return sess, ErrInValidSessId
	}
	c := p.cPool.Get()
	defer c.Close()
	dataStr, err := redis.String(c.Do("GET", p.Config.CacheKeyPrefix+sid))
	if err != nil && err != redis.ErrNil {
		return sess, ErrLoadCache
	}
	if dataStr == "" {
		return sess, ErrNoCache
	}
	err = jsoniter.UnmarshalFromString(dataStr, &sess.Values)
	if err != nil {
		return sess, ErrCacheFormat
	}
	return sess, NoErr
}

func (p *providerRedis) providerInit(conf *ProviderConf) error {
	p.providerBase.providerInit(conf)
	p.cPool = getRedisPoolByConf(conf)
	return nil
}

func (p *providerRedis) SessionDestroy(sid string) {
	key := p.Config.CacheKeyPrefix + sid
	c := p.cPool.Get()
	defer c.Close()
	c.Do("DEL", key)
}

//################# store ##########################

type sessionRedis struct {
	sessionBase
	cPool *redis.Pool
}

func (s *sessionRedis) Save() (err error) {
	if len(s.Values) == 0 {
		s.Destroy()
		return nil
	}
	if s.ID == "" {
		s.SetSessionId(GenerateUUID())
	}
	bytes, _ := jsoniter.Marshal(s.Values)
	c := s.cPool.Get()
	defer c.Close()
	key := s.Config.CacheKeyPrefix + s.ID
	_, err = redis.String(c.Do("SET", key, string(bytes), "EX", s.Config.MaxLifeTime))
	if err != nil {
		return err
	}
	s.SetCookie(s.ID)
	return nil
}

func (s *sessionRedis) SaveObject(v interface{}) (err error) {
	b, err := jsoniter.Marshal(v)
	if err != nil {
		return err
	}
	jsoniter.Unmarshal(b, &s.Values)
	return s.Save()
}

func (s *sessionRedis) Destroy() error {
	key := s.Config.CacheKeyPrefix + s.ID
	c := s.cPool.Get()
	defer c.Close()
	c.Do("DEL", key)
	s.sessionBase.Destroy()
	return nil
}
