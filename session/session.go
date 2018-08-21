package session

import (
	"sync"
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"time"
	"log"
)

type Store interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
	SessionRelease(w http.ResponseWriter)
	Flush() error
}

type Manager struct {
	cookieName string
	lock sync.Mutex
	provider Provider
	maxlifetime int64
}

func NewManager(providerName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session unknow provider %q (need to import or correct)", providerName)
	}
	return &Manager{provider:provider, cookieName:cookieName, maxlifetime:maxlifetime}, nil
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxlifeTime int64)
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

var providers = make(map[string]Provider)

func Register(name string, provider Provider)  {
	if provider == nil {
		panic("session: register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	providers[name] = provider
}

func (manager *Manager) sessionId() string{
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader,b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (manager * Manager) GetSession(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value:url.QueryEscape(sid), Path:"/", HttpOnly:true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w,&cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *Manager) DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionDestroy(cookie.Value)
	expiration := time.Now()
	cookie = &http.Cookie{Name:manager.cookieName, Path: "/", HttpOnly:true, Expires: expiration, MaxAge: -1}
	http.SetCookie(w, cookie)
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() {manager.GC()})
}

var globalSession *Manager

func SessionStart(provider, cookieName string, maxlifetime int64)  {
	var err error
	globalSession, err = NewManager(provider, cookieName, maxlifetime)
	if err != nil {
		log.Fatal(err)
	}
	go globalSession.GC()
}

func GetSession(res http.ResponseWriter, req *http.Request) Session {
	return globalSession.GetSession(res, req)
}

func DestroySession(res http.ResponseWriter, req *http.Request)  {
	globalSession.DestroySession(res, req)
}