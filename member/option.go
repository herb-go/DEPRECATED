package member

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/session"
)

//Option member service option interface
type Option interface {
	//ApplyTo  apply option to service.
	ApplyTo(*Service) error
}

//OptionFunc member service option function interface.
type OptionFunc func(*Service) error

//ApplyTo apply option function to service.
func (i OptionFunc) ApplyTo(s *Service) error {
	return i(s)
}

//OptionCommon common member service option function with  give session store.
func OptionCommon(store *session.Store) OptionFunc {
	return func(s *Service) error {
		c := cache.New()
		oc := cache.NewOptionConfig()
		oc.Driver = "dummycache"
		oc.TTL = 0
		err := oc.ApplyTo(c)
		if err != nil {
			return err
		}
		s.SessionStore = store
		s.StatusCache = c
		s.AccountsCache = c
		s.TokenCache = c
		s.RoleCache = c

		return nil
	}
}

//OptionSubCache option use sub cache node of give cache as all modules's cache.
func OptionSubCache(store *session.Store, c cache.Cacheable) OptionFunc {
	return func(s *Service) error {
		s.SessionStore = store
		s.StatusCache = cache.NewCollection(c, prefixCacheStatus, cache.DefaultTTL)
		s.AccountsCache = cache.NewCollection(c, prefixCacheAccount, cache.DefaultTTL)
		s.TokenCache = cache.NewCollection(c, prefixCacheToken, cache.DefaultTTL)
		s.RoleCache = cache.NewCollection(c, prefixCacheRole, cache.DefaultTTL)
		return nil
	}
}
