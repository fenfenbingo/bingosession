# bingosession


Golang plugin for session management (currently cookie, redis, file).You can manage your session easyly by "session.ISessionProvider" and "session.ISession" interface.

## Usage

### Start using it

Download and install it:

```bash
$ go get github.com/fenfenbingo/bingosession
```

Import it in your code:

```go
import "github.com/fenfenbingo/bingosession/session"
```

Initialize your session-provider(for once):

```
	var cookieProvider session.ISessionProvider = session.NewProvider("cookie", &session.ProviderConf{
		SessCookieName: "SESSION_COOKIE",
		CookieMaxAge:   86400 * 30,
		Path:           "/",
		Domain:         "localhost",
		Secure:         false,
		HttpOnly:       true,
		MaxLifeTime:    20,
		//specials for cookie session
		SignKey:        "qazwsx654321",
	})
```
Use your session-provider fo start or destroy a session:

```
//the interface to create or destroy session
type ISessionProvider interface {

	providerInit(conf *ProviderConf) error
    //start a session by http-request
	SessionStart(*http.Request, http.ResponseWriter) (ISession, *ErrInfo)
    //destroy a session by session_id
	SessionDestroy(sid string)
}
```
Operate your session :
```
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
	Save() error

	//extends
	LoadObject(ref interface{})
	SaveObject(v interface{}) error
}
```

## Examples

#### Cookie 

[embedmd]:# (example/cookie.go)
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
	"github.com/json-iterator/go"
)

//request http://localhost/foo to test

var cookieProvider session.ISessionProvider

func main() {
	fmt.Println("TestSessionCookie start...")
	cookieProvider = session.NewProvider("cookie", &session.ProviderConf{
		SessCookieName: "SESSION_COOKIE",
		CookieMaxAge:   86400 * 30,
		Path:           "/",
		Domain:         "localhost",
		Secure:         false,
		HttpOnly:       true,
		MaxLifeTime:    20,
		//specials for cookie session
		SignKey:        "qazwsx654321",
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", &MyHandlerCookie{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerCookie) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := cookieProvider.SessionStart(r, w)
	fmt.Println("request index:" + strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		sess.Set("uid", 555)
		sess.Set("username", "fenfen")
		sess.Save()
	} else {
		uid := jsoniter.Wrap(sess.Get("uid")).ToInt()
		usename := jsoniter.Wrap(sess.Get("username")).ToString()
		fmt.Println("get uid result:", uid)
		fmt.Println("get username result:", usename)
	}

}

type MyHandlerCookie struct {
	index int
}

```

#### Redis

[embedmd]:# (example/redis.go)
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
	"github.com/json-iterator/go"
)

//request http://localhost/foo to test

var redisProvider session.ISessionProvider

func main() {
	fmt.Println("TestSessionRedis start...")
	redisProvider = session.NewProvider("redis", &session.ProviderConf{
		SessCookieName: "SESSION_REDIS",
		CookieMaxAge:   86400 * 30,
		Path:           "/",
		Domain:         "localhost",
		Secure:         false,
		HttpOnly:       true,
		MaxLifeTime:    20,
		//specials for redis session
		CacheKeyPrefix:"sessions-",
		StoreUrl:    "127.0.0.1:6379",
		Password:    "",
		MaxIdle:     5,
		IdleTimeout: 1800,
		DBIndex:     0,
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", &MyHandlerRedis{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerRedis) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := redisProvider.SessionStart(r, w)
	fmt.Println("request index:" + strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		sess.Set("uid", 777)
		sess.Set("username", "fenfen")
		sess.Save()
	} else {
		uid := jsoniter.Wrap(sess.Get("uid")).ToInt()
		usename := jsoniter.Wrap(sess.Get("username")).ToString()
		fmt.Println("get uid result:", uid)
		fmt.Println("get username result:", usename)
	}

}

type MyHandlerRedis struct {
	index int
}

```

#### File

[embedmd]:# (example/file.go)
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
	"github.com/json-iterator/go"
)

//request http://localhost/foo to test

var fileProvider session.ISessionProvider

func main() {
	fmt.Println("TestSessionFile start...")
	fileProvider = session.NewProvider("file", &session.ProviderConf{
		SessCookieName:    "SESSION_FILE",
		CookieMaxAge:      86400 * 30,
		Path:              "/",
		Domain:            "localhost",
		Secure:            false,
		HttpOnly:          true,
		MaxLifeTime:       20,
		//specials for file session
		GCIntervalMillSec: 2000,
		StoreUrl:          "./session_cache",
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", &MyHandlerFile{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := fileProvider.SessionStart(r, w)
	fmt.Println("request index:" + strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		sess.Set("uid", 666)
		sess.Set("username", "fenfen")
		sess.Save()
	} else {
		uid := jsoniter.Wrap(sess.Get("uid")).ToInt()
		usename := jsoniter.Wrap(sess.Get("username")).ToString()
		fmt.Println("get uid result:", uid)
		fmt.Println("get username result:", usename)
	}

}

type MyHandlerFile struct {
	index int
}
```



#### Cookie&Redis (Storage by cookie,verify by Redis)

[embedmd]:# (example/cookie_redis.go)
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
	"github.com/json-iterator/go"
)

//request http://localhost/foo to test

var cookieRedisProvider session.ISessionProvider

func main() {
	fmt.Println("TestSessionCookieRedis start...")
	cookieRedisProvider = session.NewProvider("cookie_redis", &session.ProviderConf{
		SessCookieName: "SESSION_COOKIE_REDIS",
		CookieMaxAge:   86400 * 30,
		Path:           "/",
		Domain:         "localhost",
		Secure:         false,
		HttpOnly:       true,
		MaxLifeTime:    20,
		//specials for cookie session
		SignKey:        "qazwsx654321",
		CacheKeyPrefix:"sessions-",
		StoreUrl:    "127.0.0.1:6379",
		Password:    "",
		MaxIdle:     5,
		IdleTimeout: 1800,
		DBIndex:     0,
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", &MyHandlerCookieRedis{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerCookieRedis) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := cookieRedisProvider.SessionStart(r, w)
	fmt.Println("request index:" + strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		sess.Set("uid", 555)
		sess.Set("username", "fenfen")
		sess.Save()
	} else {
		uid := jsoniter.Wrap(sess.Get("uid")).ToInt()
		usename := jsoniter.Wrap(sess.Get("username")).ToString()
		fmt.Println("get uid result:", uid)
		fmt.Println("get username result:", usename)
	}

}

type MyHandlerCookieRedis struct {
	index int
}
```

#### Save custom-object (example for cookie)

[embedmd]:# (example/cookie_obj.go)
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
)

//request http://localhost/foo to test

var cookieObjProvider session.ISessionProvider

func main() {
	fmt.Println("TestSessionCookieObj start...")
	cookieObjProvider = session.NewProvider("cookie", &session.ProviderConf{
		SessCookieName:    "SESSION_COOKIE_OBJ",
		CookieMaxAge:      86400 * 30,
		Path:              "/",
		Domain:            "localhost",
		Secure:            false,
		HttpOnly:          true,
		MaxLifeTime:       20,
		//specials for cookie session
		SignKey:"qazwsx654321",
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", &MyHandlerCookieObj{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerCookieObj) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := cookieObjProvider.SessionStart(r, w)
	fmt.Println("request index:"+strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		user:=&TestUserInfo{Uid:666,Username:"fenfen"}
		sess.SaveObject(user)
	} else {
		user:=&TestUserInfo{}
		sess.LoadObject(user)
		fmt.Println("get uid result:", user.Uid)
		fmt.Println("get username result:", user.Username)
	}

}

type MyHandlerCookieObj struct {
	index int
}

type TestUserInfo struct {
	Uid int64 `json:"uid"`
	Username string `json:"username"`
}
```

#### Destroy a session when no user request (not support for cookie-only session,example for file)

[embedmd]:# (example/destroy.go)
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
	"github.com/json-iterator/go"
	"time"
)

//request http://localhost/foo to test

var fileProvider2 session.ISessionProvider

func main() {
	fmt.Println("TestSessionFile(Destroy) start...")
	fileProvider2 = session.NewProvider("file", &session.ProviderConf{
		SessCookieName: "SESSION_FILE",
		CookieMaxAge:   86400 * 30,
		Path:           "/",
		Domain:         "localhost",
		Secure:         false,
		HttpOnly:       true,
		MaxLifeTime:    86400,
		//specials for file session
		GCIntervalMillSec: 2000,
		StoreUrl:          "./session_cache",
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", &MyHandlerFile2{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerFile2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := fileProvider2.SessionStart(r, w)
	fmt.Println("request index:" + strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		sess.Set("uid", 666)
		sess.Set("username", "fenfen")
		sess.Save()
		sid := sess.GetSessionId()
		go func() {
			//Although the session expires in 86400 seconds, due to active destruction, the session will expire in 5 seconds.
			//Use "ISessionProvider.SessionDestroy" function,you can expire old sessions when you modify your password
			time.Sleep(time.Second * 5)
			fileProvider2.SessionDestroy(sid)
		}()
	} else {
		uid := jsoniter.Wrap(sess.Get("uid")).ToInt()
		usename := jsoniter.Wrap(sess.Get("username")).ToString()
		fmt.Println("get uid result:", uid)
		fmt.Println("get username result:", usename)
	}

}

type MyHandlerFile2 struct {
	index int
}

```