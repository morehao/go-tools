package esquery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildQuery(t *testing.T) {
	boolQuery := NewBool().
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
			).BuildQuery(),
		)

	dsl := NewSearchBody().
		WithQuery(boolQuery.BuildQuery()).
		SortBy("balance", false).
		FromVal(0).
		SizeVal(10).
		AggsMap("average_balance", AggAvg("balance")).
		Build()

	queryStr, err := dsl.String()
	assert.Nil(t, err)

	t.Log(queryStr)
}
