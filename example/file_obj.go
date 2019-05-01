package main

import (
	"fmt"
	"net/http"
	"github.com/fenfenbingo/bingosession/session"
	"strconv"
)

//request http://localhost/foo to test

var fileObjProvider session.ISessionProvider

func main() {
	fmt.Println("TestSessionFileObj start...")
	fileObjProvider = session.NewProvider("file", &session.ProviderConf{
		SessCookieName:    "SESSION_FILE_OBJ",
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
	mux.Handle("/foo", &MyHandlerFileObj{})
	http.ListenAndServe(":80", mux)
}

func (h *MyHandlerFileObj) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, ecode := fileObjProvider.SessionStart(r, w)
	fmt.Println("request index:" + strconv.Itoa(h.index))
	h.index++
	if ecode != session.NoErr {
		fmt.Println("session start fail,err_code:", ecode)
		user:=&TestUserInfo2{Uid:666,Username:"fenfen"}
		sess.SaveObject(user)
	} else {
		user:=&TestUserInfo2{}
		sess.LoadObject(user)
		fmt.Println("get uid result:", user.Uid)
		fmt.Println("get username result:", user.Username)
	}

}

type MyHandlerFileObj struct {
	index int
}

type TestUserInfo2 struct {
	Uid int64 `json:"uid"`
	Username string `json:"username"`
}