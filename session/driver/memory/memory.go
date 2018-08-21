package memory

import (
	"sync"
	"container/list"
	"time"
	"github.com/pengzj/cherry/session"
)

type Provider struct {
	lock sync.Mutex
	sessions map[string]*list.Element
	list *list.List
}


type SessionStore struct {
	sid string
	timeAccessed time.Time
	values map[interface{}]interface{}
}

var pder = &Provider{list: list.New()}


func (st *SessionStore) Set(key, value interface{}) error {
	st.values[key] = value
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.sid)
	if v, ok := st.values[key];ok {
		return v
	}
	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.values, key)
	pder.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

func (provider *Provider) SessionInit(sid string) (session.Session, error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	v := make(map[interface{}]interface{},0)
	newsess := &SessionStore{sid:sid, timeAccessed:time.Now(), values: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (provider *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid];ok {
		return element.Value.(*SessionStore), nil
	}

	session, err := pder.SessionInit(sid)
	return session, err
}

func (provider *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
	}
	return nil
}

func (provider *Provider) SessionGC(maxlifetime int64)  {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (provider *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
	}
	return nil
}

func init()  {
	pder.sessions = make(map[string]*list.Element,0)
	session.Register("memory", pder)
}


