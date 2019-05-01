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
