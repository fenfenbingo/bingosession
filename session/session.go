package session

import (
	"net/http"
	"github.com/fenfenbingo/bingosession/internal/util"
)

type ProviderConf struct {
	//configs for http-headers
	SessCookieName string
	CookieMaxAge   int
	Path           string
	Domain         string
	Secure         bool
	HttpOnly       bool

	//缓存生存时间（秒）
	//cache life time for seconds
	MaxLifeTime int64

	//定时清理无效session的时间间隔(file类型)
	GCIntervalMillSec int64

	//SignKey用于cookie存储的session签名，改变SignKey会使所有session失效
	//SignKey is  used in cookie session,it will  make all sessions of this provider expire if you change it
	SignKey string

	StoreUrl string
	//CacheKeyPrefix用于redis存储中key的前缀，改变CacheKeyPrefix会使所有session失效
	//CacheKeyPrefix is  used in redis session,it will make all sessions of this provider expire if you change it
	CacheKeyPrefix string

	//configs for redis
	Password    string
	DBIndex     int
	MaxIdle     int
	IdleTimeout int64
}

//the interface to operate session
type ISession interface {
	GetSessionId() string

	SetSessionId(string)

	Get(key string) interface{}

	Set(key string, val interface{})

	Delete(key string)

	Destroy() error

	//如果本来SessionId为空，调用Save()后会生成一个32位全局唯一ID，可以通过GetSessionId()方法获取。也可以使用SetSessionId()事先设置SessionId。
	//生成SessionId请使用GenerateUUID()方法，不支持自定义格式。
	//If the SessionId was empty,after "Save()" ,a 32-bit UUID is generated,get it by called "GetSessionId()".
	//You can use "GenerateUUID()" to generate in advance too,and use "SetSessionId()" to set it.
	Save() error

	//extends(扩展方法)
	LoadObject(ref interface{})
	SaveObject(v interface{}) error
}

//the interface to start or destroy session
type ISessionProvider interface {
	providerInit(conf *ProviderConf) error
	//SessionDestroy用于用户请求时,根据请求内容获取session缓存
	SessionStart(*http.Request, http.ResponseWriter) (ISession, *ErrInfo)
	//SessionDestroy用于非用户请求时,根据session_id清除session缓存
	SessionDestroy(sid string)
}

func NewProvider(cacheType string, c *ProviderConf) ISessionProvider {
	var res ISessionProvider
	switch cacheType {
	case "file":
		{
			res = &providerFile{}
		}
	case "cookie":
		{
			res = &providerCookie{}
		}
	case "cookie_redis":
		{
			res = &providerCookieRedis{}
		}
	case "redis":
		{
			res = &providerRedis{}
		}
	}
	err := res.providerInit(c)
	if err != nil {
		return nil
	}
	return res
}

//生成32位全局唯一ID
func GenerateUUID() string {
	return util.GenerateUUID()
}
