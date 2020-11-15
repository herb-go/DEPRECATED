package member

import (
	"strconv"
	"time"

	"github.com/herb-go/deprecated/cache/datastore"
	"github.com/herb-go/user"
)

type testAccountProvider struct {
	AccountsMap map[string]user.Accounts
}

func newTestAccount(uid string) *user.Account {
	account, err := service.NewAccount("test", uid)
	if err != nil {
		panic(err)
	}
	return account
}
func (s *testAccountProvider) Execute(service *Service) error {
	service.AccountsProvider = s
	return nil
}
func (s *testAccountProvider) Accounts(uid ...string) (*Accounts, error) {
	var a = Accounts{}
	for _, v := range uid {
		a[v] = s.AccountsMap[v]
	}
	return &a, nil
}
func (s *testAccountProvider) AccountToUID(account *user.Account) (uid string, err error) {
	for uid, v := range s.AccountsMap {
		if v.Exists(account) {
			return uid, nil
		}
	}
	return "", nil
}
func (s *testAccountProvider) Register(account *user.Account) (uid string, err error) {
	for _, v := range s.AccountsMap {
		if v.Exists(account) {
			return "", ErrAccountRegisterExists
		}
	}
	uid = strconv.Itoa(len(s.AccountsMap))
	s.AccountsMap[uid] = []*user.Account{account}
	return uid, nil
}
func (s *testAccountProvider) AccountToUIDOrRegister(account *user.Account) (uid string, registerd bool, err error) {
	for uid, v := range s.AccountsMap {
		if v.Exists(account) {
			return uid, false, nil
		}
	}
	s.AccountsMap[account.Account] = []*user.Account{account}
	return account.Account, true, nil
}
func (s *testAccountProvider) BindAccount(uid string, account *user.Account) error {
	for _, v := range s.AccountsMap {
		if v.Exists(account) {
			return user.ErrAccountBindingExists
		}
	}
	if s.AccountsMap[uid] == nil {
		s.AccountsMap[uid] = user.Accounts{}
	}
	accounts := s.AccountsMap[uid]
	err := accounts.Bind(account)
	if err != nil {
		return err
	}
	s.AccountsMap[uid] = accounts
	return nil
}
func (s *testAccountProvider) UnbindAccount(uid string, account *user.Account) error {
	if s.AccountsMap[uid] == nil {
		return user.ErrAccountUnbindingNotExists
	}
	accounts := s.AccountsMap[uid]
	err := accounts.Unbind(account)
	if err != nil {
		return err
	}
	s.AccountsMap[uid] = accounts
	return nil
}

func newTestAccountProvider() *testAccountProvider {
	return &testAccountProvider{
		AccountsMap: map[string]user.Accounts{},
	}
}

type testTokenProvider struct {
	MemberTokens map[string]string
}

func (s *testTokenProvider) Execute(service *Service) error {
	service.TokenProvider = s
	return nil
}
func (s *testTokenProvider) Tokens(uid ...string) (Tokens, error) {
	var r = Tokens{}
	for _, v := range uid {
		r[v] = s.MemberTokens[v]
	}
	return r, nil

}
func (s *testTokenProvider) Revoke(uid string) (string, error) {
	var ts = strconv.FormatInt(time.Now().UnixNano(), 10)
	s.MemberTokens[uid] = ts
	return ts, nil
}
func newTestTokenProvider() *testTokenProvider {
	return &testTokenProvider{
		MemberTokens: map[string]string{},
	}
}

type testStatusService struct {
	StatusMap StatusMap
}

func (s *testStatusService) Execute(service *Service) error {
	service.StatusProvider = s
	return nil
}

func (s *testStatusService) Statuses(uid ...string) (StatusMap, error) {
	var r = StatusMap{}
	for _, v := range uid {
		r[v] = s.StatusMap[v]
	}
	return r, nil

}
func (s *testStatusService) SetStatus(uid string, status Status) error {
	v := status
	s.StatusMap[uid] = v
	return nil
}

//SupportedStatus return supported status map
func (s *testStatusService) SupportedStatus() map[Status]bool {
	return StatusMapAll
}

func newTestStatusProvider() *testStatusService {
	return &testStatusService{
		StatusMap: StatusMap{},
	}
}

type testPasswordProvider struct {
	Passwords map[string]string
}

func (p *testPasswordProvider) PasswordChangeable() bool {
	return false
}

func (s *testPasswordProvider) Execute(service *Service) error {
	service.PasswordProvider = s
	return nil
}
func (s *testPasswordProvider) VerifyPassword(uid string, password string) (bool, error) {

	pass := s.Passwords[uid]
	if pass == "" {
		return false, nil
	}
	return pass == password, nil

}
func (s *testPasswordProvider) UpdatePassword(uid string, password string) error {
	s.Passwords[uid] = password
	return nil
}

func newTestPasswordProvider() *testPasswordProvider {
	return &testPasswordProvider{
		Passwords: map[string]string{},
	}
}

type userProfiles map[string]map[string][]string

var rawUserProfiles = map[string]map[string][]string{}

func (p userProfiles) NewMapElement(key string) error {
	p[key] = map[string][]string{}
	return nil
}
func (p userProfiles) LoadMapElements(keys ...string) error {
	for _, v := range keys {
		p[v] = rawUserProfiles[v]
	}
	return nil
}
func (p userProfiles) Map() interface{} {
	return p
}
func newTestUesrProfiles() *userProfiles {
	return &userProfiles{}
}

var testUserProfileProvider = &datastore.DataSource{
	Creator: func() interface{} {
		return map[string][]string{}
	},
	SourceLoader: func(keys ...string) (map[string]interface{}, error) {
		result := map[string]interface{}{}
		for _, v := range keys {
			result[v] = rawUserProfiles[v]
		}
		return result, nil
	},
}

type testRoleProvider Roles

func (s *testRoleProvider) Execute(service *Service) error {
	service.RoleProvider = s
	return nil
}
func (s *testRoleProvider) Roles(uid ...string) (*Roles, error) {
	r := Roles{}
	for _, v := range uid {
		r[v] = (*s)[v]
	}
	return &r, nil

}

func newTestRoleProvider() *testRoleProvider {
	return &testRoleProvider{}
}
