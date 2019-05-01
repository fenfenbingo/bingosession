package session

import (
	"net/http"
	"io/ioutil"
	"github.com/json-iterator/go"
	"path"
	"sync"
	"time"
	"os"
	"path/filepath"
)

var filelock sync.Mutex

//session_redis 使用文件存放数据

//############### provider ######################

type providerFile struct {
	providerBase
}

func (p *providerFile) SessionStart(req *http.Request, w http.ResponseWriter) (res ISession, ecode *ErrInfo) {
	sid := p.GetCookieContent(req)
	sess := &sessionFile{
		sessionBase{
			Values:     make(map[string]interface{}),
			Req:        req,
			RespWriter: w,
			ID:         sid,
			Config:     p.Config,
		},
	}
	if sid == "" {
		return sess, ErrEmpty
	}
	defer func() {
		if ecode != NoErr {
			res.SetSessionId("")
		}
	}()
	if len(sid) != 32 {
		return sess, ErrInValidSessId
	}
	filelock.Lock()
	bytes, err := ioutil.ReadFile(path.Join(p.StoreUrl, sid[0:2], sid))
	filelock.Unlock()
	if err != nil {
		return sess, ErrLoadCache
	}
	err = jsoniter.Unmarshal(bytes, &sess.Values)
	if err != nil {
		return sess, ErrCacheFormat
	}
	return sess, NoErr
}

// remove expired files
func (p *providerFile) gcfiles(path string, info os.FileInfo, err error) (er error) {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if (info.ModTime().Unix() + p.Config.MaxLifeTime) < time.Now().Unix() {
		filelock.Lock()
		defer filelock.Unlock()
		os.Remove(path)
	}
	return nil
}

func (p *providerFile) providerInit(conf *ProviderConf) (err error) {
	p.StoreUrl = conf.StoreUrl
	p.Config = conf
	go p.GCCall()
	return nil
}

func (p *providerFile) GCCall() {
	time.Sleep(time.Duration(p.Config.GCIntervalMillSec) * time.Millisecond)
	filepath.Walk(p.Config.StoreUrl, p.gcfiles)
	p.GCCall()
	return
}

func (p *providerFile) SessionDestroy(sid string) {
	filelock.Lock()
	defer filelock.Unlock()
	os.Remove(path.Join(p.Config.StoreUrl, string(sid[0:2]), sid))
}

//############ store ####################

type sessionFile struct {
	sessionBase
}

func (s *sessionFile) Save() (err error) {
	if len(s.Values) == 0 {
		s.Destroy()
		return nil
	}
	if s.ID == "" {
		s.SetSessionId(GenerateUUID())
	}
	return s.save(s.Values)
}

func (s *sessionFile) save(v interface{}) (err error) {
	bytes, err := jsoniter.Marshal(v)
	if err != nil {
		return err
	}
	filedir := path.Join(s.Config.StoreUrl, string(s.ID[0:2]))
	filepath := path.Join(s.Config.StoreUrl, string(s.ID[0:2]), s.ID)
	filelock.Lock()
	defer filelock.Unlock()
	_, err = os.Stat(filedir)
	if err != nil {
		os.MkdirAll(filedir, 0777)
	}
	err = ioutil.WriteFile(filepath, bytes, 0777)
	if err != nil {
		return
	}
	s.SetCookie(s.ID)
	return nil
}

func (s *sessionFile) SaveObject(v interface{}) (err error) {
	if s.ID == "" {
		s.SetSessionId(GenerateUUID())
	}
	return s.save(v)
}

func (s *sessionFile) Destroy() (err error) {
	if s.ID != "" {
		filelock.Lock()
		os.Remove(path.Join(s.Config.StoreUrl, string(s.ID[0:2]), s.ID))
		filelock.Unlock()
	}
	s.sessionBase.Destroy()
	return nil
}
