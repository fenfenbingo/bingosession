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
