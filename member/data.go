package member

//DEPRECATED
import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/datastore"
)

//ServiceData member user data module.
//DEPRECATED
type ServiceData struct {
	service *Service
}

//Cache Return member user data cache.
//DEPRECATED
func (s *ServiceData) Cache(field string) cache.Cacheable {
	return s.service.DataCache
}

//Clean clean member user data cache by uid.
//DEPRECATED
func (s *ServiceData) Clean(field string, uid string) error {
	return s.Cache(field).Del(uid)
}

//Load load and cache user data map from provider.
//Return any error if raised.
//DEPRECATED
func (s *ServiceData) Load(field string, data datastore.Store, keys ...string) error {
	p := s.service.DataProviders[field]
	return datastore.Load(
		data,
		s.Cache(field),
		p.SourceLoader,
		p.Creator,
		keys...,
	)
}
