package uidaccount

import (
	"errors"
	"strings"

	"github.com/herb-go/user"
	"github.com/herb-go/deprecated/member"
)

// ErrUIDAndAccountNotMatch error raised when uid and account not match
var ErrUIDAndAccountNotMatch = errors.New("uidaccount:uid and account not match")

// ErrAccountKeywordNotMatch error raised when account keyword not match
var ErrAccountKeywordNotMatch = errors.New("uidaccount: account keyword not match")

//ErrPrefixOrSuffixNotMatch error raised when prefix or suffix not match
var ErrPrefixOrSuffixNotMatch = errors.New("uidaccount:prefix or suffix not match")

//UIDAccount uidaccount directive struct
type UIDAccount struct {
	//AccountKeyword account keyword
	AccountKeyword string
	//Prefix uid prefix
	Prefix string
	//Suffix uid suffix
	Suffix string
}

func (u *UIDAccount) accountToUID(account string) (string, error) {
	if strings.HasPrefix(account, u.Prefix) && strings.HasSuffix(account, u.Suffix) {
		return strings.TrimSuffix(strings.TrimPrefix(account, u.Prefix), u.Suffix), nil
	}
	return "", ErrPrefixOrSuffixNotMatch
}
func (u *UIDAccount) uidToAccount(uid string) string {
	return u.Prefix + uid + u.Suffix
}

//Accounts return account map of given uid list.
//Return account map and any error if raised.
func (u *UIDAccount) Accounts(uid ...string) (*member.Accounts, error) {
	a := member.Accounts{}
	for _, id := range uid {
		if id == "" {
			continue
		}
		account := user.NewAccount()
		account.Keyword = u.AccountKeyword
		account.Account = u.uidToAccount(id)
		a[id] = user.Accounts{account}
	}
	return &a, nil
}

//AccountToUID query uid by user account.
//Return user id and any error if raised.
//Return empty string as userid if account not found.
func (u *UIDAccount) AccountToUID(account *user.Account) (uid string, err error) {
	if account.Keyword != u.AccountKeyword {
		return "", ErrAccountKeywordNotMatch
	}
	return u.accountToUID(account.Account)
}

//Register create new user with given account.
//Return created user id and any error if raised.
//Privoder should return ErrAccountRegisterExists if account is used.
func (u *UIDAccount) Register(account *user.Account) (uid string, err error) {
	return u.AccountToUID(account)
}

//AccountToUIDOrRegister query uid by user account.Register user if account not found.
//Return user id and any error if raised.
func (u *UIDAccount) AccountToUIDOrRegister(account *user.Account) (uid string, registerd bool, err error) {
	uid, err = u.AccountToUID(account)
	return uid, false, err
}

//BindAccount bind account to user.
//Return any error if raised.
//If account exists,user.ErrAccountBindingExists should be rasied.
func (u *UIDAccount) BindAccount(uid string, account *user.Account) error {
	if account.Keyword != u.AccountKeyword {
		return ErrAccountKeywordNotMatch
	}
	id, err := u.accountToUID(account.Account)
	if err != nil {
		return err
	}
	if uid != id {
		return ErrUIDAndAccountNotMatch
	}
	return nil
}

//UnbindAccount unbind account from user.
//Return any error if raised.
//If account not exists,user.ErrAccountUnbindingNotExists should be rasied.
func (u *UIDAccount) UnbindAccount(uid string, account *user.Account) error {
	if account.Keyword != u.AccountKeyword {
		return ErrAccountKeywordNotMatch
	}
	id, err := u.accountToUID(account.Account)
	if err != nil {
		return err
	}
	if uid != id {
		return ErrUIDAndAccountNotMatch
	}
	return nil

}

//Execute apply uidaccount directive to member service
func (u *UIDAccount) Execute(m *member.Service) error {
	m.AccountsProvider = u
	return nil
}

//DirectiveFactory factory to create uidaccount directive
var DirectiveFactory = func(loader func(v interface{}) error) (member.Directive, error) {
	c := &UIDAccount{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
