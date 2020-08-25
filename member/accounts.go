package member

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/datastore"
	"github.com/herb-go/user"
)

//Accounts user accounts map.
type Accounts map[string]user.Accounts

//AccountsStore accounts data store
type AccountsStore struct {
	*datastore.SyncMapStore
}

//Get get acounts by given user id
func (s *AccountsStore) Get(uid string) user.Accounts {
	v, ok := s.Load(uid)
	if !ok {
		return nil
	}
	if v == nil {
		return nil
	}
	return *(v.(*user.Accounts))
}

//NewAccountsStore create new account data store
func NewAccountsStore() *AccountsStore {
	return &AccountsStore{
		datastore.NewSyncMapStore(),
	}
}

//ServiceAccounts Member accounts module.
type ServiceAccounts struct {
	service *Service
}

func (s *ServiceAccounts) loader(keys ...string) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := s.service.AccountsProvider.Accounts(keys...)
	if err != nil {
		return result, err
	}
	result = map[string]interface{}{}
	for k := range *data {
		v := (*data)[k]
		result[k] = &v
	}
	return result, nil
}

//Cache Return member accounts cache.
func (s *ServiceAccounts) Cache() cache.Cacheable {
	return s.service.AccountsCache
}

//Clean clean accounts cache by uid.
//Return any error if raised.
func (s *ServiceAccounts) Clean(uid string) error {
	return s.Cache().Del(uid)
}

//Load load and cache accounts from provider.
//Return any error if raised.
func (s *ServiceAccounts) Load(accounts datastore.Store, keys ...string) error {
	return datastore.Load(
		accounts,
		s.Cache(),
		s.loader,
		func() interface{} {
			return &user.Accounts{}
		},
		keys...,
	)
}

//Register create new user with given account.
//Return created user id and any error if raised.
func (s *ServiceAccounts) Register(account *user.Account) (uid string, err error) {
	return s.service.AccountsProvider.Register(account)
}

//AccountToUID query uid by user account.
//Return user id and any error if raised.
//Return empty string as userid if account not found.
func (s *ServiceAccounts) AccountToUID(account *user.Account) (uid string, err error) {
	return s.service.AccountsProvider.AccountToUID(account)
}

//AccountToUIDOrRegister query uid by user account.Register user if account not found.
//Return user id ,whether registered and any error if raised.
func (s *ServiceAccounts) AccountToUIDOrRegister(account *user.Account) (uid string, registerd bool, err error) {
	return s.service.AccountsProvider.AccountToUIDOrRegister(account)
}

//BindAccount bind account to user.
//user account cache will be cleand.
//Return any error if raised.
//If account exists,user.ErrAccountBindingExists should be rasied.
func (s *ServiceAccounts) BindAccount(uid string, account *user.Account) error {
	err := s.service.AccountsProvider.BindAccount(uid, account)
	if err != nil {
		return err
	}
	return s.Clean(uid)
}

//UnbindAccount unbind account from user.
//user account cache will be cleand.
//Return any error if raised.
//If account not exists,user.ErrAccountUnbindingNotExists should be rasied.
func (s *ServiceAccounts) UnbindAccount(uid string, account *user.Account) error {
	err := s.service.AccountsProvider.UnbindAccount(uid, account)
	if err != nil {
		return err
	}
	return s.Clean(uid)
}

//AccountsProvider member account provider interface
type AccountsProvider interface {
	//Accounts return account map of given uid list.
	//Return account map and any error if raised.
	Accounts(uid ...string) (*Accounts, error)
	//AccountToUID query uid by user account.
	//Return user id and any error if raised.
	//Return empty string as userid if account not found.
	AccountToUID(account *user.Account) (uid string, err error)
	//Register create new user with given account.
	//Return created user id and any error if raised.
	//Privoder should return ErrAccountRegisterExists if account is used.
	Register(account *user.Account) (uid string, err error)
	//AccountToUIDOrRegister query uid by user account.Register user if account not found.
	//Return user id and any error if raised.
	AccountToUIDOrRegister(account *user.Account) (uid string, registerd bool, err error)
	//BindAccount bind account to user.
	//Return any error if raised.
	//If account exists,user.ErrAccountBindingExists should be rasied.
	BindAccount(uid string, account *user.Account) error
	//UnbindAccount unbind account from user.
	//Return any error if raised.
	//If account not exists,user.ErrAccountUnbindingNotExists should be rasied.
	UnbindAccount(uid string, account *user.Account) error
}
