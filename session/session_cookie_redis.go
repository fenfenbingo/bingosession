package session

import (
	"github.com/gomodule/redigo/redis"
	"net/http"
)

//session_cookie_redis 使用cookie存放数据，仅使用redis验证是否过期

//########### provider ##############

type providerCookieRedis struct {
	providerCookie
	cPool *redis.Pool
}

func (p *providerCookieRedis) providerInit(conf *ProviderConf) error{
	p.providerCookie.providerInit(conf)
	p.cPool = getRedisPoolByConf(conf)
	return nil
}

func (p *providerCookieRedis) SessionStart(req *http.Request, w http.ResponseWriter) (res ISession, ecode *ErrInfo) {
	s, ecode := p.providerCookie.sessionStart(req, w)
	res = &sessionCookieRedis{
		*s,
		p.cPool,
	}
	if ecode != NoErr {
		return
	}
	c := p.cPool.Get()
	defer c.Close()
	key := p.Config.CacheKeyPrefix + s.ID
	if ok, _ := redis.Bool(c.Do("EXISTS", key)); !ok {
		return res, ErrExpire
	}
	return
}

func (p *providerCookieRedis) SessionDestroy(sid string) {
	key := p.Config.CacheKeyPrefix + sid
	c := p.cPool.Get()
	defer c.Close()
	c.Do("DEL", key)
}

//########### store ##############

type sessionCookieRedis struct {
	sessionCookie
	cPool *redis.Pool
}

func (s *sessionCookieRedis) Save() (err error) {
	err = s.sessionCookie.Save()
	if err != nil {
		return
	}
	key := s.Config.CacheKeyPrefix + s.ID
	c:=s.cPool.Get()
	defer c.Close()
	_, err = c.Do("SET", key, "1", "EX", s.Config.MaxLifeTime)
	return
}

func (s *sessionCookieRedis) Destroy() error {
	key := s.Config.CacheKeyPrefix + s.ID
	s.cPool.Get().Do("DEL", key)
	s.sessionCookie.Destroy()
	return nil
}
