package member

import (
	"github.com/herb-go/deprecated/cache/datastore"
	"github.com/herb-go/user/profile"
)

//Profiles  user profiles map type
type Profiles map[string]*profile.Profile

func (p *Profiles) Chain(profiles *Profiles) {
	for k, v := range *profiles {
		profile, ok := (*p)[k]
		if ok {
			profile.Chain(v)
		} else {
			(*p)[k] = v.Clone()
		}
	}
}

//ProfilesStore user profiles data store
type ProfilesStore struct {
	*datastore.SyncMapStore
}

// Get get user profiles by given user id
func (s *ProfilesStore) Get(uid string) *profile.Profile {
	v, ok := s.Load(uid)
	if !ok {
		return nil
	}
	if v == nil {
		return nil
	}
	return v.(*profile.Profile)
}

//NewProfilesStore create new user profiles data store
func NewProfilesStore() *ProfilesStore {
	return &ProfilesStore{
		datastore.NewSyncMapStore(),
	}
}

//ServiceProfiles member profile module.
type ServiceProfiles struct {
	service *Service
}

//Load load user profiles from provider.
//Return any error if raised.
func (s *ServiceProfiles) Load(store datastore.Store, keys ...string) error {
	var result = Profiles{}
	for _, v := range s.service.ProfilesProviders {
		p, err := v.Profiles(keys...)
		if err != nil {
			return err
		}
		result.Chain(p)
	}
	for k, v := range result {
		store.Store(k, v)
	}
	return nil
}
func (s *ServiceProfiles) UpdateProfile(uid string, profile *profile.Profile) error {
	for _, v := range s.service.ProfilesProviders {
		err := v.UpdateProfile(uid, profile)
		if err != nil {
			return err
		}
	}
	return nil
}

//ProfilesProvider  member role provider interface
type ProfilesProvider interface {
	Profiles(uid ...string) (*Profiles, error)
	UpdateProfile(uid string, profile *profile.Profile) error
}
