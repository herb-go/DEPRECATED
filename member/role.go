package member

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/datastore"
	"github.com/herb-go/herbsecurity/authorize/role"
)

//Roles  user role map type
type Roles map[string]*role.Roles

//RolesStore user roles data store
type RolesStore struct {
	*datastore.SyncMapStore
}

// Get get user roles by given user id
func (s *RolesStore) Get(uid string) *role.Roles {
	v, ok := s.Load(uid)
	if !ok {
		return nil
	}
	if v == nil {
		return nil
	}
	return v.(*role.Roles)
}

//NewRolesStore create new user rolees data store
func NewRolesStore() *RolesStore {
	return &RolesStore{
		datastore.NewSyncMapStore(),
	}
}

//ServiceRole member role module.
type ServiceRole struct {
	service *Service
}

//Load load and cache user roles from provider.
//Return any error if raised.
func (s *ServiceRole) Load(store datastore.Store, keys ...string) error {
	return datastore.Load(
		store,
		s.Cache(),
		s.loader,
		func() interface{} {
			return &role.Roles{}
		},
		keys...,
	)
}

//Cache Return member role cache.
func (s *ServiceRole) Cache() cache.Cacheable {
	return s.service.RoleCache
}

//Clean clean role cache by uid.
func (s *ServiceRole) Clean(uid string) error {
	return s.Cache().Del(uid)
}
func (s *ServiceRole) loader(keys ...string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	data, err := s.service.RoleProvider.Roles(keys...)
	if err != nil {
		return result, err
	}
	for k := range *data {
		v := (*data)[k]
		result[k] = v
	}
	return result, nil
}

//RolesProvider  member role provider interface
type RolesProvider interface {
	Roles(uid ...string) (*Roles, error)
}
