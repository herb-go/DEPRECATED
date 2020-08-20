package member

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/datastore"
)

//Tokens user token map type
type Tokens map[string]string

// TokensStore user token data store
type TokensStore struct {
	*datastore.SyncMapStore
}

// Get get user token by given user id
func (s *TokensStore) Get(uid string) string {
	v, ok := s.Load(uid)
	if !ok {
		return ""
	}
	if v == nil {
		return ""
	}
	return *(v.(*string))
}

// NewTokensStore create new user token data store.
func NewTokensStore() *TokensStore {
	return &TokensStore{
		datastore.NewSyncMapStore(),
	}
}

//ServiceToken member token module.
type ServiceToken struct {
	service *Service
}

//TokenProvider  member token provider interface
type TokenProvider interface {
	Tokens(uid ...string) (Tokens, error)
	Revoke(uid string) (string, error)
}

//Cache Return member token cache.
func (s *ServiceToken) Cache() cache.Cacheable {
	return s.service.TokenCache
}

//Clean clean token cache by uid.
func (s *ServiceToken) Clean(uid string) error {
	return s.Cache().Del(uid)
}

//Revoke revoke user token and regenerate new token.
//user revoke cache will be cleand.
//Return new token and any error if resied.
func (s *ServiceToken) Revoke(uid string) (string, error) {
	t, err := s.service.TokenProvider.Revoke(uid)
	if err != nil {
		return "", err
	}
	err = s.Clean(uid)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (s *ServiceToken) loader(keys ...string) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := s.service.TokenProvider.Tokens(keys...)
	if err != nil {
		return result, err
	}
	result = map[string]interface{}{}
	for k := range data {
		v := data[k]
		result[k] = &v
	}
	return result, nil

}

//Load load and cache token from provider.
//Return any error if raised.
func (s *ServiceToken) Load(Tokens datastore.Store, keys ...string) error {
	return datastore.Load(
		Tokens,
		s.Cache(),
		s.loader,
		func() interface{} {
			var s = ""
			return &s
		},
		keys...,
	)
}
