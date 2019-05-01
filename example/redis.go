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
