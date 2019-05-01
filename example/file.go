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
