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
