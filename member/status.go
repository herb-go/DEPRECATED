package member

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/datastore"
)

//StatusNormal user status normal
const StatusNormal = Status(0)

// StatusBanned user status banned
const StatusBanned = Status(1)

// StatusRevoked user status revoked
const StatusRevoked = Status(2)

// StatusPending user status pending
const StatusPending = Status(3)

// StatusExpired user status expried
const StatusExpired = Status(4)

var StatusMapAll = map[Status]bool{
	StatusNormal:  true,
	StatusBanned:  true,
	StatusRevoked: true,
	StatusPending: true,
	StatusExpired: true,
}

var StatusMapMin = map[Status]bool{
	StatusNormal: true,
	StatusBanned: true,
}

//Status user status type
type Status int

// IsAvaliable  check if user status is normal status.
func (s *Status) IsAvaliable() bool {
	return IsAvaliable(s)
}

//StatusMap user  status map
//User is banned if if map data of user id is true
type StatusMap map[string]Status

//StatusStore user status data store
type StatusStore struct {
	*datastore.SyncMapStore
}

// IsAvaliable  check if user status is normal status.
func IsAvaliable(s *Status) bool {
	if s == nil {
		return false
	}
	return *s == StatusNormal
}

// Get get user status by given user id
func (s *StatusStore) Get(uid string) *Status {
	v, ok := s.Load(uid)
	if !ok {
		return nil
	}
	if v == nil {
		return nil
	}
	return v.(*Status)
}

// NewStatusStore create new status data store
func NewStatusStore() *StatusStore {
	return &StatusStore{
		datastore.NewSyncMapStore(),
	}
}

//StatusProvider  member  status provider interface
type StatusProvider interface {
	//Statuses return  status  map of given uid list.
	//Return status  map and any error if raised.
	Statuses(uid ...string) (StatusMap, error)
	//SetStatus set user status.
	//Return any error if raised.
	SetStatus(uid string, status Status) error
	//SupportedStatus return supported status map
	SupportedStatus() map[Status]bool
}

//ServiceStatus Member  status module.
type ServiceStatus struct {
	service *Service
}

//Load load and cache user  status from provider.
//Return any error if raised.
func (s *ServiceStatus) Load(statusMap datastore.Store, keys ...string) error {
	return datastore.Load(
		statusMap,
		s.Cache(),
		s.loader,
		func() interface{} {
			var val = StatusNormal
			return &val
		},
		keys...,
	)
}

//Cache Return member  status cache.
func (s *ServiceStatus) Cache() cache.Cacheable {
	return s.service.StatusCache
}

//Clean clean  status cache by uid.
func (s *ServiceStatus) Clean(uid string) error {
	return s.Cache().Del(uid)
}

//SetStatus set user  status.
//user status cache will be cleand.
//Return any error if raised.
func (s *ServiceStatus) SetStatus(uid string, status Status) error {
	ok := s.service.StatusProvider.SupportedStatus()[status]
	if !ok {
		return ErrStatusNotSupport
	}
	err := s.service.StatusProvider.SetStatus(uid, status)
	if err != nil {
		return err
	}
	return s.Clean(uid)

}

func (s *ServiceStatus) loader(keys ...string) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := s.service.StatusProvider.Statuses(keys...)
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
