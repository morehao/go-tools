package dbes

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Status int    `json:"status"`
}

func TestDSLBuilder(t *testing.T) {
	// 创建ES客户端
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// 创建文档
	user := &User{
		ID:     "1",
		Name:   "张三",
		Age:    25,
		Status: 1,
	}
	resp, err := NewBuilder(client).
		SetIndex("users").
		SetID(user.ID).
		SetDoc(user).
		Create(ctx)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	t.Logf("Create Response Status: %s", resp.Status())

	// 更新文档
	updateResp, err := NewBuilder(client).
		SetIndex("users").
		SetID("1").
		SetDoc(M(
			"doc", M(
				"age", 26,
				"status", 2,
			),
		)).
		Update(ctx)

	if err != nil {
		t.Fatal(err)
	}
	defer updateResp.Body.Close()
	t.Logf("Update Response Status: %s", updateResp.Status())

	// 根据查询条件批量更新
	updateByQueryResp, err := NewBuilder(client).
		SetIndex("users").
		SetQuery(Q.Term("status", 1)).
		Set("script", M(
			"source", "ctx._source.status = params.status",
			"lang", "painless",
			"params", M(
				"status", 2,
			),
		)).
		UpdateByQuery(ctx)

	if err != nil {
		t.Fatal(err)
	}
	defer updateByQueryResp.Body.Close()
	t.Logf("Update By Query Response Status: %s", updateByQueryResp.Status())

	// 删除文档
	deleteResp, err := NewBuilder(client).
		SetIndex("users").
		SetID("1").
		Delete(ctx)

	if err != nil {
		t.Fatal(err)
	}
	defer deleteResp.Body.Close()
	t.Logf("Delete Response Status: %s", deleteResp.Status())

	// 根据查询条件批量删除
	deleteByQueryResp, err := NewBuilder(client).
		SetIndex("users").
		SetQuery(Q.Term("status", 2)).
		DeleteByQuery(ctx)

	if err != nil {
		t.Fatal(err)
	}
	defer deleteByQueryResp.Body.Close()
	t.Logf("Delete By Query Response Status: %s", deleteByQueryResp.Status())

	// 批量操作
	bulk := NewBulkBuilder().
		// 批量创建
		AddCreate("users", "2", &User{
			Name:   "李四",
			Age:    30,
			Status: 1,
		}).
		AddCreate("users", "3", &User{
			Name:   "王五",
			Age:    35,
			Status: 1,
		}).
		// 批量更新
		AddUpdate("users", "2", M(
			"age", 31,
			"status", 2,
		)).
		// 批量删除
		AddDelete("users", "3")

	bulkResp, err := NewBuilder(client).
		SetIndex("users").
		Bulk(ctx, bulk.Build())

	if err != nil {
		t.Fatal(err)
	}
	defer bulkResp.Body.Close()
	t.Logf("Bulk Response Status: %s", bulkResp.Status())

	// 搜索示例
	searchResp, err := NewBuilder(client).
		SetIndex("users").
		SetQuery(Q.Bool(
			[]map[string]interface{}{
				Q.Match("name", "李四"),
			},
			nil,
			nil,
			[]map[string]interface{}{
				Q.Range("age", M(
					"gte", 30,
					"lt", 40,
				)),
			},
		)).
		SetSort([]interface{}{
			S.Field("age", "desc"),
		}).
		SetSize(10).
		Search(ctx)

	if err != nil {
		t.Fatal(err)
	}
	defer searchResp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(searchResp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	t.Logf("Search Result: %+v", result)
}
