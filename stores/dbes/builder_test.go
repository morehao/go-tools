package dbes

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
	ctx := context.Background()

	logCfg := &glog.ModuleLoggerConfig{
		Service:   "ES",
		Level:     glog.InfoLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"requestId"},
	}
	opt := glog.WithZapOptions(zap.AddCallerSkip(2))
	initLogErr := glog.NewLogger(logCfg, opt)
	assert.Nil(t, initLogErr)
	cfg := ESConfig{
		Service: "es",
		Addr:    "http://localhost:9200",
	}
	client, _, initErr := InitES(cfg)
	assert.Nil(t, initErr)
	defer glog.Close()

	t.Run("CreateAccount", func(t *testing.T) {
		// 查询当前最大的 ID
		account := &Account{
			AccountNumber: 1,
			Balance:       10000,
			Firstname:     "Alice",
			Lastname:      "Smith",
			Age:           30,
			Gender:        "F",
			Email:         "alice@example.com",
			Employer:      "TechCorp",
			Address:       "123 Main St",
			City:          "New York",
			State:         "NY",
		}
		resp, err := NewBuilder(client).
			SetIndex("accounts").
			SetID("1").
			SetDoc(account).
			Create(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		t.Logf("Create Response Status: %s", resp.Status())
	})

	t.Run("UpdateBalance", func(t *testing.T) {
		resp, err := NewBuilder(client).
			SetIndex("accounts").
			SetID("1").
			SetDoc(M("doc", M("balance", 15000))).
			Update(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		t.Logf("Update Balance Status: %s", resp.Status())
	})

	t.Run("SearchByCity", func(t *testing.T) {
		resp, err := NewBuilder(client).
			SetIndex("accounts").
			SetQuery(Q.Match("city", "Brogan")).
			SetSize(5).
			Search(ctx)
		assert.Nil(t, err)
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}
		t.Logf("Search Result: %s", gutils.ToJsonString(result))
	})

	t.Run("BulkInsert", func(t *testing.T) {
		bulk := NewBulkBuilder().
			AddCreate("accounts", "2", &Account{
				AccountNumber: 2,
				Balance:       8000,
				Firstname:     "Bob",
				Lastname:      "Brown",
				City:          "Chicago",
				State:         "IL",
			}).
			AddCreate("accounts", "3", &Account{
				AccountNumber: 3,
				Balance:       20000,
				Firstname:     "Charlie",
				Lastname:      "White",
				City:          "Los Angeles",
				State:         "CA",
			})
		resp, err := NewBuilder(client).
			SetIndex("accounts").
			Bulk(ctx, bulk.Build())
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		t.Logf("Bulk Insert Status: %s", resp.Status())
	})

	t.Run("DeleteAccount", func(t *testing.T) {
		resp, err := NewBuilder(client).
			SetIndex("accounts").
			SetID("1").
			Delete(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		t.Logf("Delete Account Status: %s", resp.Status())
	})
}
