package dbes

import (
	"testing"

	"github.com/morehao/go-tools/glog"
)

type Account struct {
	AccountNumber int64  `json:"account_number"`
	Balance       int64  `json:"balance"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Age           int    `json:"age"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	Employer      string `json:"employer"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
}

func TestAccountsIndex(t *testing.T) {
	query := NewBuilder().
		SetIndex("accounts").
		SetQuery(Q.Match("city", "Brogan")).
		SetSize(5).Build()
	t.Logf("Search Result: %s", glog.ToJsonString(query))
}
