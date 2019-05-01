package session

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func getRedisPoolByConf(conf *ProviderConf) (p *redis.Pool) {
	p = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.StoreUrl)
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				if _, err = c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if conf.DBIndex > 0 {
				_, err = c.Do("SELECT", conf.DBIndex)
				if err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
	}
	return p
}
