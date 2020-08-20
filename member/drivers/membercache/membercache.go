package membercache

import (
	"errors"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/deprecated/member"
)

//CacheTypeAll cache type for user cache as all type.
var CacheTypeAll = ""

//CacheTypeStatus cache type for use as status cache only.
var CacheTypeStatus = "status"

//CacheTypeAccounts cache type for use as accounts cache only.
var CacheTypeAccounts = "accounts"

// CacheTypeToken cache type for use as token cache only
var CacheTypeToken = "token"

// CacheTypeRole cache type for use as role cache only
var CacheTypeRole = "role"

// CacheTypeData cache type for use as data cache only
var CacheTypeData = "data"

//ErrUnknownMemberCacheType error raised when cache type unknown.
var ErrUnknownMemberCacheType = errors.New("membercache:unknown member cache type")

//Config member cache config struct
type Config struct {
	Cache *cache.OptionConfig
	Type  string
}

// Execute apply config to member service
func (c *Config) Execute(m *member.Service) error {
	membercache := cache.New()
	err := c.Cache.ApplyTo(membercache)
	if err != nil {
		return err
	}
	switch c.Type {
	case CacheTypeAll:
		m.StatusCache = cache.NewCollection(membercache, "Status", cache.DefaultTTL)
		m.AccountsCache = cache.NewCollection(membercache, "Account", cache.DefaultTTL)
		m.TokenCache = cache.NewCollection(membercache, "Token", cache.DefaultTTL)
		m.RoleCache = cache.NewCollection(membercache, "Role", cache.DefaultTTL)
		m.DataCache = cache.NewNode(membercache, "data")
		return nil
	case CacheTypeStatus:
		m.StatusCache = membercache
		return nil
	case CacheTypeToken:
		m.TokenCache = membercache
		return nil
	case CacheTypeAccounts:
		m.AccountsCache = membercache
		return nil
	case CacheTypeRole:
		m.RoleCache = membercache
		return nil
	case CacheTypeData:
		m.DataCache = membercache
		return nil

	}
	return ErrUnknownMemberCacheType
}

//DirectiveFactory factory to create cache directive
var DirectiveFactory = func(loader func(v interface{}) error) (member.Directive, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
