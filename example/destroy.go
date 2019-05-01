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
