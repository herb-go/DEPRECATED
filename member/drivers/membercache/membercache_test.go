package membercache_test

import (
	"testing"

	"github.com/herb-go/herbconfig/loader"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"
	"github.com/herb-go/deprecated/member"
	"github.com/herb-go/deprecated/member/drivers/membercache"
)

type DirectiveConfig struct {
	Config func(v interface{}) error `config:", lazyload"`
}

func TestMemberCacheAll(t *testing.T) {
	dummycache := cache.Dummy()
	m := member.New()
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
	config := &DirectiveConfig{}
	err := loader.LoadConfig("json", []byte(configAll), config)
	if err != nil {
		panic(err)
	}
	d, err := membercache.DirectiveFactory(config.Config)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != nil {
		panic(err)
	}
	if m.AccountsCache == dummycache ||
		m.RoleCache == dummycache ||
		m.StatusCache == dummycache ||
		m.TokenCache == dummycache ||
		m.DataCache == dummycache {
		t.Fatal(m)
	}
}

func TestMemberCacheStatus(t *testing.T) {
	dummycache := cache.Dummy()
	m := member.New()
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
	config := &DirectiveConfig{}
	err := loader.LoadConfig("json", []byte(configStatus), config)
	if err != nil {
		panic(err)
	}
	d, err := membercache.DirectiveFactory(config.Config)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != nil {
		panic(err)
	}
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache == dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
}

func TestMemberCacheAccounts(t *testing.T) {
	dummycache := cache.Dummy()
	m := member.New()
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
	config := &DirectiveConfig{}
	err := loader.LoadConfig("json", []byte(configAccounts), config)
	if err != nil {
		panic(err)
	}
	d, err := membercache.DirectiveFactory(config.Config)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != nil {
		panic(err)
	}
	if m.AccountsCache == dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
}

func TestMemberCacheToken(t *testing.T) {
	dummycache := cache.Dummy()
	m := member.New()
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
	config := &DirectiveConfig{}
	err := loader.LoadConfig("json", []byte(configToken), config)
	if err != nil {
		panic(err)
	}
	d, err := membercache.DirectiveFactory(config.Config)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != nil {
		panic(err)
	}
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache == dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
}

func TestMemberCacheRole(t *testing.T) {
	dummycache := cache.Dummy()
	m := member.New()
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
	config := &DirectiveConfig{}
	err := loader.LoadConfig("json", []byte(configRole), config)
	if err != nil {
		panic(err)
	}
	d, err := membercache.DirectiveFactory(config.Config)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != nil {
		panic(err)
	}
	if m.AccountsCache != dummycache ||
		m.RoleCache == dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
}

func TestMemberCacheData(t *testing.T) {
	dummycache := cache.Dummy()
	m := member.New()
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache != dummycache {
		t.Fatal(m)
	}
	config := &DirectiveConfig{}
	err := loader.LoadConfig("json", []byte(configData), config)
	if err != nil {
		panic(err)
	}
	d, err := membercache.DirectiveFactory(config.Config)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != nil {
		panic(err)
	}
	if m.AccountsCache != dummycache ||
		m.RoleCache != dummycache ||
		m.StatusCache != dummycache ||
		m.TokenCache != dummycache ||
		m.DataCache == dummycache {
		t.Fatal(m)
	}
}

func TestMemberCacheError(t *testing.T) {
	m := member.New()
	config := &DirectiveConfig{}
	err := loader.LoadConfig("json", []byte(configError), config)
	if err != nil {
		panic(err)
	}
	d, err := membercache.DirectiveFactory(config.Config)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != membercache.ErrUnknownMemberCacheType {
		t.Fatal(err)
	}
}
