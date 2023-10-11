package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ai-feier/geek-web/session"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	errorKeyNotFound     = errors.New("session: 找不到 key")
	errorSessionNotFound = errors.New("session: 找不到 session")
)

type Store struct {
	mutex      sync.RWMutex
	// 将 Session 存入缓存
	sessions   *cache.Cache
	// session 默认过期时间
	expiration time.Duration
}

type StoreOpts func(store *Store)

func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessions:   cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess := &Session{
		id: id,
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.sessions.Get(id)
	if !ok {
		return fmt.Errorf("session: 该 id 对应的 session 不存在 %s", id)
	}
	s.sessions.Set(id, val, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.sessions.Get(id)
	if !ok {
		return nil, errorSessionNotFound
	}
	return val.(*Session), nil
}

type Session struct {
	// sessionId
	id string

	// mutex sync.RWMutex
	// values map[string]any
	// session 实际存储类型
	values sync.Map
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, ok := s.values.Load(key)
	if !ok {
		// return nil, fmt.Errorf("%w, key %s", errorKeyNotFound, key)
		return nil, errorKeyNotFound
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	s.values.Store(key, val)
	return nil
}

func (s *Session) ID() string {
	return s.id
}
