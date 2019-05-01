package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
)

//request http://localhost/foo to test

var redisObjProvider session.ISessionProvider

func main() {
	fmt.Println("TestSessionRedisObj start...")
	redisObjProvider = session.NewProvider("redis", &session.ProviderConf{
		SessCookieName:    "SESSION_REDIS_OBJ",
		CookieMaxAge:      86400 * 30,
		Path:              "/",
		Domain:            "localhost",
		Secure:            false,
		HttpOnly:          true,
		MaxLifeTime:       20,
		//specials for redis session
		CacheKeyPrefix:"sessions-",
		StoreUrl:    "127.0.0.1:6379",
		Password:    "",
		MaxIdle:     5,
		IdleTimeout: 1800,
		DBIndex:     0,
	})
	mux := http.NewServeMux()
	mux.Handle("/foo", &MyHandlerRedisObj{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerRedisObj) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := redisObjProvider.SessionStart(r, w)
	fmt.Println("request index:"+strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		user:=&TestUserInfo1{Uid:666,Username:"fenfen"}
		sess.SaveObject(user)
	} else {
		user:=&TestUserInfo1{}
		sess.LoadObject(user)
		fmt.Println("get uid result:", user.Uid)
		fmt.Println("get username result:", user.Username)
	}

}

type MyHandlerRedisObj struct {
	index int
}

type TestUserInfo1 struct {
	Uid int64 `json:"uid"`
	Username string `json:"username"`
}

