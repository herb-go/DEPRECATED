package uidaccount

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/herb-go/user"

	"github.com/herb-go/deprecated/member"
)

func TestUIDAccount(t *testing.T) {
	var testUIDAccount = &UIDAccount{
		AccountKeyword: "testkeyword",
		Prefix:         "testprefix",
		Suffix:         "testsuffix",
	}
	data, err := json.Marshal(testUIDAccount)
	if err != nil {
		panic(err)
	}
	loader := json.NewDecoder(bytes.NewBuffer(data))

	m := member.New()
	d, err := DirectiveFactory(loader.Decode)
	if err != nil {
		panic(err)
	}
	err = d.Execute(m)
	if err != nil {
		panic(err)
	}
	as := member.NewAccountsStore()
	err = m.Accounts().Load(as, "test", "")
	accounts := as.Get("test")
	if len(accounts) != 1 || accounts[0].Keyword != testUIDAccount.AccountKeyword || accounts[0].Account != testUIDAccount.Prefix+"test"+testUIDAccount.Suffix {
		t.Fatal(as)
	}
	uid, err := m.Accounts().AccountToUID(accounts[0])
	if err != nil {
		panic(err)
	}
	if uid != "test" {
		t.Fatal(m)
	}
	wrongaccount := user.NewAccount()
	wrongaccount.Keyword = testUIDAccount.AccountKeyword
	wrongaccount.Account = "test"
	uid, err = m.Accounts().AccountToUID(wrongaccount)
	if err != ErrPrefixOrSuffixNotMatch || uid != "" {
		t.Fatal(uid, err)
	}
	wrongaccount = user.NewAccount()
	wrongaccount.Keyword = "wrongkeyword"
	wrongaccount.Account = testUIDAccount.Prefix + "test" + testUIDAccount.Suffix
	uid, err = m.Accounts().AccountToUID(wrongaccount)
	if err != ErrAccountKeywordNotMatch || uid != "" {
		t.Fatal(uid, err)
	}
	uid, err = m.Accounts().Register(accounts[0])
	if err != nil {
		panic(err)
	}
	if uid != "test" {
		t.Fatal(m)
	}
	uid, registered, err := m.Accounts().AccountToUIDOrRegister(accounts[0])
	if err != nil {
		panic(err)
	}
	if uid != "test" || registered != false {
		t.Fatal(m)
	}
	err = m.Accounts().BindAccount("test", accounts[0])
	if err != nil {
		panic(err)
	}
	wrongaccount = user.NewAccount()
	wrongaccount.Keyword = testUIDAccount.AccountKeyword
	wrongaccount.Account = "test"
	err = m.Accounts().BindAccount("test", wrongaccount)
	if err != ErrPrefixOrSuffixNotMatch {
		panic(err)
	}
	wrongaccount = user.NewAccount()
	wrongaccount.Keyword = "wrongkeyword"
	wrongaccount.Account = testUIDAccount.Prefix + "test" + testUIDAccount.Suffix
	err = m.Accounts().BindAccount("test", wrongaccount)
	if err != ErrAccountKeywordNotMatch {
		panic(err)
	}
	wrongaccount = user.NewAccount()
	wrongaccount.Keyword = testUIDAccount.AccountKeyword
	wrongaccount.Account = testUIDAccount.Prefix + "wrongtest" + testUIDAccount.Suffix
	err = m.Accounts().BindAccount("test", wrongaccount)
	if err != ErrUIDAndAccountNotMatch {
		panic(err)
	}

	err = m.Accounts().UnbindAccount("test", accounts[0])
	if err != nil {
		panic(err)
	}

	wrongaccount = user.NewAccount()
	wrongaccount.Keyword = testUIDAccount.AccountKeyword
	wrongaccount.Account = "test"
	err = m.Accounts().UnbindAccount("test", wrongaccount)
	if err != ErrPrefixOrSuffixNotMatch {
		panic(err)
	}
	wrongaccount = user.NewAccount()
	wrongaccount.Keyword = "wrongkeyword"
	wrongaccount.Account = testUIDAccount.Prefix + "test" + testUIDAccount.Suffix
	err = m.Accounts().UnbindAccount("test", wrongaccount)
	if err != ErrAccountKeywordNotMatch {
		panic(err)
	}
	wrongaccount = user.NewAccount()
	wrongaccount.Keyword = testUIDAccount.AccountKeyword
	wrongaccount.Account = testUIDAccount.Prefix + "wrongtest" + testUIDAccount.Suffix
	err = m.Accounts().UnbindAccount("test", wrongaccount)
	if err != ErrUIDAndAccountNotMatch {
		panic(err)
	}
}
