package esquery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComplexNestedQuery(t *testing.T) {
	// 构建主 bool 查询
	boolQuery := NewBool().
		Filter(Term("state.keyword", "IL")). // 确保使用 keyword 字段
		Must(
			Range("age", Query{"gt": 30}),
			Range("balance", Query{"gte": 20000, "lte": 50000}),
		).
		Should(
			Wildcard("lastname", "*son*"),
			NewBool().Must(
				Term("employer", "Scentric"),
				Match("email", "ratliff"),
			).BuildQuery(), // 嵌套 Bool 子查询
		)

	// 构建 SearchBody
	body := NewSearchBody(boolQuery.BuildQuery()).
		SortBy("balance", false).
		SetFrom(0).
		SetSize(10).
		SetAgg("average_balance", AggAvg("balance")).
		SetAgg("group_by_state", AggTerms("state.keyword", 5))

	// 转成 JSON 字符串
	queryStr, err := body.ToBuffer()
	assert.NoError(t, err)
	t.Log(string(queryStr.Bytes()))
}
