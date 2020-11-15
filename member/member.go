package member

import (
	"github.com/herb-go/deprecated/cache/datastore"
)

//Members members stand for cached user data.
type Members struct {
	Service     *Service
	Accounts    *AccountsStore
	StatusStore *StatusStore
	Tokens      *TokensStore
	Roles       *RolesStore
	Profiles    *ProfilesStore
	Dataset     map[string]datastore.Store
}

//LoadStatus load banned status for users.
//loaded status will stored in members StatusStore field.
//Return  status map and any error if rased.
func (m *Members) LoadStatus(keys ...string) (*StatusStore, error) {
	return m.StatusStore, m.Service.Status().Load(m.StatusStore, keys...)
}

//Status return user status.
//Return user  status and any error if rased.
func (m *Members) Status(key string) (*Status, error) {
	smap, err := m.LoadStatus(key)
	if err != nil {
		return nil, err
	}
	return smap.Get(key), nil
}

//LoadTokens load tokens for users.
//loaded tokens will stored in members Tokens field.
//Return Tokens and any error if rased.
func (m *Members) LoadTokens(keys ...string) (*TokensStore, error) {
	return m.Tokens, m.Service.Token().Load(m.Tokens, keys...)
}

//LoadAccount load accounts for users.
//loaded tokens will stored in members Accounts field.
//Return Accounts and any error if rased.
func (m *Members) LoadAccount(keys ...string) (*AccountsStore, error) {
	return m.Accounts, m.Service.Accounts().Load(m.Accounts, keys...)
}

//LoadRoles load roles for users.
//loaded roles will stored in members Roles field.
//Return Roles and any error if rased.
func (m *Members) LoadRoles(keys ...string) (*RolesStore, error) {
	return m.Roles, m.Service.Roles().Load(m.Roles, keys...)
}

//LoadProfiles load profiles for users.
//loaded profiles will stored in members Profiles field.
//Return Profiles and any error if rased.
func (m *Members) LoadProfiles(keys ...string) (*ProfilesStore, error) {
	return m.Profiles, m.Service.Profiles().Load(m.Roles, keys...)
}

//LoadData load named data for users.
//loaded datas will stored in members Dataset field.
//Return datas and any error if rased.
func (m *Members) LoadData(field string, keys ...string) (datastore.Store, error) {
	return m.Dataset[field], m.Service.Data().Load(field, m.Dataset[field], keys...)
}

//Data return named data field of members
func (m *Members) Data(field string) datastore.Store {
	return m.Dataset[field]
}

//NewMembers return empty members with given service.
func NewMembers(s *Service) *Members {
	var member = &Members{
		Service:     s,
		Accounts:    NewAccountsStore(),
		StatusStore: NewStatusStore(),
		Roles:       NewRolesStore(),
		Tokens:      NewTokensStore(),
		Profiles:    NewProfilesStore(),
	}
	member.Dataset = make(map[string]datastore.Store, len(s.DataProviders))
	for k := range s.DataProviders {
		member.Dataset[k] = datastore.NewSyncMapStore()
	}
	return member
}
