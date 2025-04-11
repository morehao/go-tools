package esquery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildQuery(t *testing.T) {
	query := NewBool().
		Filter(Term("state", "IL")).
		Must(
			Range("age", map[string]any{"gt": 30}),
			Range("balance", map[string]any{"gte": 20000, "lte": 50000}),
		).
		Should(
			Wildcard("lastname", "*son*"),
			NewBool().Must(
				Term("employer", "Scentric"),
				Match("email", "ratliff"),
			).BuildQuery(), // 返回嵌套查询
		).
		Sort("balance", false).
		From(0).
		Size(10).
		Aggs("average_balance", AggAvg("balance"))

	// 调用 Build 获取查询结果（结构化查询对象）
	queryResult := query.BuildDSL()

	// 将查询对象转换为 JSON 字符串
	queryStr, err := queryResult.String()
	assert.Nil(t, err)

	// 打印查询结果
	t.Log(queryStr)
}
