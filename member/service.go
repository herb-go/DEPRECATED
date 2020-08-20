package member

import (
	"net/http"

	"github.com/herb-go/herb/user/profile"

	"github.com/herb-go/protecter"

	"github.com/herb-go/herbsecurity/authorize/role"

	"context"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/datastore"
	"github.com/herb-go/herb/user"
	"github.com/herb-go/herb/user/httpuser"
	"github.com/herb-go/session"
)

const prefixCacheStatus = "S"
const prefixCacheAccount = "A"
const prefixCacheToken = "T"
const prefixCacheRole = "R"

//DefaultSessionUIDFieldName default user id session field name when create member service.
const DefaultSessionUIDFieldName = "herb-member-uid"

//DefaultSessionMemberTokenFieldName default member token session field name when create member service.
const DefaultSessionMemberTokenFieldName = "herb-member-token"

//ContextType member context name type.
type ContextType string

//DefaultContextName default member context name.
var DefaultContextName = ContextType("members")

//Service member service main interafce.
type Service struct {
	//SessionStore session store which stores member data.
	SessionStore *session.Store
	//SessionUIDFieldName session field which stores user id.
	SessionUIDFieldName string
	//SessionMemberFieldName session field which stores member token.
	SessionMemberFieldName string
	//ContextName context name stores members data.
	ContextName ContextType
	//BannedProvider user banned status provider.
	//DON'T use this provider directly,use Service.Banned() instead.
	StatusProvider StatusProvider
	//BannedCache data stores banned status.
	StatusCache cache.Cacheable
	//AccountsProvider user accounts provider.
	//DON'T use this provider directly,use Service.Accounts() instead.
	AccountsProvider AccountsProvider
	//AccountsCache data stores user accounts.
	AccountsCache cache.Cacheable
	//TokenProvider user token provider.
	//DON'T use this provider directly,use Service.Tokens() instead.
	TokenProvider TokenProvider
	//TokenCache data stores user tokens.
	TokenCache cache.Cacheable
	//PasswordProvider user password provider.
	//DON'T use this provider directly,use Service.Password() instead.
	PasswordProvider PasswordProvider
	//RoleProvider user roles provider.
	//DON'T use this provider directly,use Service.Roles() instead.
	RoleProvider RolesProvider
	//RoleCache data stores user roles.
	RoleCache cache.Cacheable
	//DataProviders user data provider.
	//A map of registered data map type.
	DataProviders map[string]*datastore.DataSource
	//DataCache data stores user dataset.
	DataCache cache.Cacheable
	//User Profiles Providers
	ProfilesProviders []ProfilesProvider
	//AccountProviders registered account provider map.
	AccountProviders map[string]user.AccountProvider
}

func (s *Service) Reset() {
	s.SessionStore = nil
	s.SessionUIDFieldName = ""
	s.SessionMemberFieldName = ""
	s.ContextName = ""
	s.StatusProvider = nil
	s.AccountsProvider = nil
	s.TokenProvider = nil
	s.PasswordProvider = nil
	s.RoleProvider = nil
	s.RoleProvider = nil
	s.DataProviders = map[string]*datastore.DataSource{}
	s.AccountProviders = map[string]user.AccountProvider{}
	s.StatusCache = cache.Dummy()
	s.AccountsCache = cache.Dummy()
	s.TokenCache = cache.Dummy()
	s.RoleCache = cache.Dummy()
	s.DataCache = cache.Dummy()

}

//RegisterAccountProvider register account provider as keyword.
func (s *Service) RegisterAccountProvider(keyword string, t user.AccountProvider) {
	s.AccountProviders[keyword] = t
}

//NewAccount create new account by given registered keyword and account name.
//Return created user account and any error if raised.
//Return ErrAccountKeywordNotRegistered if account keyword is not registered by Service.RegisterAccountType .
func (s *Service) NewAccount(keyword string, account string) (*user.Account, error) {
	accountType, ok := s.AccountProviders[keyword]
	if ok == false {
		return nil, ErrAccountKeywordNotRegistered
	}
	return accountType.NewAccount(keyword, account)
}

//Accounts return Accounts module.
func (s *Service) Accounts() *ServiceAccounts {
	return &ServiceAccounts{
		service: s,
	}
}

//Password return Password modules.
func (s *Service) Password() *ServicePassword {
	return &ServicePassword{
		service: s,
	}
}

//Status return Status modules.
func (s *Service) Status() *ServiceStatus {
	return &ServiceStatus{
		service: s,
	}
}

//Token return Token modules.
func (s *Service) Token() *ServiceToken {
	return &ServiceToken{
		service: s,
	}
}

//Data return Data modules.
//DEPRECATED
func (s *Service) Data() *ServiceData {
	return &ServiceData{
		service: s,
	}
}

//Roles return Roles modules.
func (s *Service) Roles() *ServiceRole {
	return &ServiceRole{
		service: s,
	}
}

//Profiles return profiles modules.
func (s *Service) Profiles() *ServiceProfiles {
	return &ServiceProfiles{
		service: s,
	}
}

//RegisterData register data type as named data field.
//data type should implement DataProvider interface so that data module can create and load user data.
//Return any error if raised.
//DEPRECATED
func (s *Service) RegisterData(key string, p *datastore.DataSource) error {
	s.DataProviders[key] = p
	return nil
}

//NewMembers return new members data.
func (s *Service) NewMembers() *Members {
	return NewMembers(s)
}

//GetMembersFromRequest get members data in http request context.
//Create new members data and bind to context if not exist.
//Return members data.
func (s *Service) GetMembersFromRequest(r *http.Request) (members *Members) {
	var contextName = s.ContextName
	if contextName == "" {
		contextName = DefaultContextName
	}
	var membersInterface = r.Context().Value(contextName)
	if membersInterface != nil {
		if members, ok := membersInterface.(*Members); ok == true {
			return members
		}
	}
	members = NewMembers(s)
	var ctx = context.WithValue(r.Context(), contextName, members)
	*r = *r.WithContext(ctx)
	return members
}

//UIDField return user id session field
func (s *Service) UIDField() *session.Field {
	var fieldName = s.SessionUIDFieldName
	if fieldName == "" {
		fieldName = DefaultSessionUIDFieldName
	}
	return s.SessionStore.Field(fieldName)
}

//MemberTokenField return member  token session field
func (s *Service) MemberTokenField() *session.Field {
	var fieldName = s.SessionMemberFieldName
	if fieldName == "" {
		fieldName = DefaultSessionMemberTokenFieldName
	}
	return s.SessionStore.Field(fieldName)
}

//IdentifyRequest Identify user in http request.
//Return user id and any error raised.
//If user is not logged in,returned user id will by empty string.
func (s *Service) IdentifyRequest(r *http.Request) (uid string, err error) {
	uid, err = s.UIDField().IdentifyRequest(r)
	if err != nil {
		return "", err
	}
	if uid == "" {
		return "", nil
	}
	var members = s.GetMembersFromRequest(r)
	if s.TokenProvider != nil {
		_, err = members.LoadTokens(uid)
		if err != nil {
			return "", err
		}
		var token string
		err = s.MemberTokenField().Get(r, &token)
		if err != nil {
			return "", err
		}
		if token != members.Tokens.Get(uid) {
			return "", nil
		}
	}
	return uid, nil
}

func (s *Service) RequestProfiles(r *http.Request) (*profile.Profile, error) {
	uid, err := s.IdentifyRequest(r)
	if err != nil {
		return nil, err
	}
	var members = s.GetMembersFromRequest(r)
	_, err = members.LoadProfiles(uid)
	if err != nil {
		return nil, err
	}
	return members.Profiles.Get(uid), nil
}

//Logout Logout user in http request.
func (s *Service) Logout(w http.ResponseWriter, r *http.Request) error {
	return s.UIDField().Logout(w, r)
}

//Authorizer create Authorizer with given rule provider.
func (s *Service) Authorizer(rs protecter.PolicyLoader) httpuser.Authorizer {
	return &Authorizer{
		Service:      s,
		PolicyLoader: rs,
	}
}

//LoginRequiredMiddleware return login requred middleware with given unauthorizedAction.
func (s *Service) LoginRequiredMiddleware(unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return httpuser.LoginRequiredMiddleware(s, unauthorizedAction)
}

//LogoutMiddleware return logout middleware.
func (s *Service) LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return httpuser.LogoutMiddleware(s)
}

//AuthorizeMiddleware return Authorize Middleware with special rule provider.
//Middleware will check user banned status if banned status provider is installed.
func (s *Service) AuthorizeMiddleware(rs protecter.PolicyLoader, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return httpuser.AuthorizeMiddleware(s.Authorizer(rs), unauthorizedAction)
}

//BannedMiddleware return Authorize Middleware check only if user is banned
func (s *Service) BannedMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return httpuser.AuthorizeMiddleware(s.Authorizer(nil), nil)
}

//RolesAuthorizeMiddleware return Authorize Middleware with roles as rule provider.
//Middleware will check user banned status if banned status provider is installed.
func (s *Service) RolesAuthorizeMiddleware(ruleNames ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rs = role.NewPlainRoles(ruleNames...)
	return s.AuthorizeMiddleware(protecter.RolePolicyLoader(rs), nil)
}

//Login login giver user to http request
func (s *Service) Login(w http.ResponseWriter, r *http.Request, id string) error {
	err := s.UIDField().Login(w, r, id)
	if err != nil {
		return err
	}
	if s.TokenProvider != nil {
		member := s.GetMembersFromRequest(r)
		tokens, err := member.LoadTokens(id)
		if err != nil {
			return err
		}
		err = s.MemberTokenField().Set(r, tokens.Get(id))
		if err != nil {
			return err
		}
	}
	return nil
}

//Init servcei with given option.
func (s *Service) Init(option Option) error {
	return option.ApplyTo(s)
}

//New create new member service with given session store.
func New() *Service {
	return &Service{
		DataProviders:    map[string]*datastore.DataSource{},
		AccountProviders: map[string]user.AccountProvider{},
		StatusCache:      cache.Dummy(),
		AccountsCache:    cache.Dummy(),
		TokenCache:       cache.Dummy(),
		RoleCache:        cache.Dummy(),
		DataCache:        cache.Dummy(),
	}
}
