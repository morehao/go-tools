package esquery

import jsoniter "github.com/json-iterator/go"

type Query map[string]any

type BoolQuery struct {
	MustClauses    []Query
	ShouldClauses  []Query
	MustNotClauses []Query
	FilterClauses  []Query
}

func NewBool() *BoolQuery {
	return &BoolQuery{}
}

func (b *BoolQuery) Must(q ...Query) *BoolQuery {
	b.MustClauses = append(b.MustClauses, q...)
	return b
}

func (b *BoolQuery) Should(q ...Query) *BoolQuery {
	b.ShouldClauses = append(b.ShouldClauses, q...)
	return b
}

func (b *BoolQuery) MustNot(q ...Query) *BoolQuery {
	b.MustNotClauses = append(b.MustNotClauses, q...)
	return b
}

func (b *BoolQuery) Filter(q ...Query) *BoolQuery {
	b.FilterClauses = append(b.FilterClauses, q...)
	return b
}

func (b *BoolQuery) BuildQuery() Query {
	boolQuery := Query{}
	if len(b.MustClauses) > 0 {
		boolQuery["must"] = b.MustClauses
	}
	if len(b.ShouldClauses) > 0 {
		boolQuery["should"] = b.ShouldClauses
	}
	if len(b.MustNotClauses) > 0 {
		boolQuery["must_not"] = b.MustNotClauses
	}
	if len(b.FilterClauses) > 0 {
		boolQuery["filter"] = b.FilterClauses
	}
	return Query{"bool": boolQuery}
}

type SearchBody struct {
	Query   Query
	Sort    []Query
	From    *int
	Size    *int
	Aggs    Query
	Options []func(Query)
}

func NewSearchBody() *SearchBody {
	return &SearchBody{}
}

func (b *SearchBody) WithQuery(q Query) *SearchBody {
	b.Query = q
	return b
}

func (b *SearchBody) SortBy(field string, asc bool) *SearchBody {
	order := "asc"
	if !asc {
		order = "desc"
	}
	b.Sort = append(b.Sort, Query{field: Query{"order": order}})
	return b
}

func (b *SearchBody) FromVal(from int) *SearchBody {
	b.From = &from
	return b
}

func (b *SearchBody) SizeVal(size int) *SearchBody {
	b.Size = &size
	return b
}

func (b *SearchBody) AggsMap(name string, agg Query) *SearchBody {
	if b.Aggs == nil {
		b.Aggs = Query{}
	}
	b.Aggs[name] = agg
	return b
}

func (b *SearchBody) With(opt func(Query)) *SearchBody {
	b.Options = append(b.Options, opt)
	return b
}

func (b *SearchBody) Build() Query {
	result := Query{
		"query": b.Query,
	}
	if len(b.Sort) > 0 {
		result["sort"] = b.Sort
	}
	if b.From != nil {
		result["from"] = *b.From
	}
	if b.Size != nil {
		result["size"] = *b.Size
	}
	if len(b.Aggs) > 0 {
		result["aggs"] = b.Aggs
	}
	for _, opt := range b.Options {
		opt(result)
	}
	return result
}

func (q Query) String() (string, error) {
	return jsoniter.MarshalToString(q)
}

func Term(field string, value any) Query {
	return Query{"term": Query{field: Query{"value": value}}}
}

func Match(field string, value any) Query {
	return Query{"match": Query{field: value}}
}

func Range(field string, op Query) Query {
	return Query{"range": Query{field: op}}
}

func Exists(field string) Query {
	return Query{"exists": Query{"field": field}}
}

func Wildcard(field, pattern string) Query {
	return Query{"wildcard": Query{field: Query{"value": pattern}}}
}

func Prefix(field, prefix string) Query {
	return Query{"prefix": Query{field: prefix}}
}

func Script(script string) Query {
	return Query{"script": Query{"source": script}}
}

func Nested(path string, query Query) Query {
	return Query{"nested": Query{
		"path":  path,
		"query": query,
	}}
}

func Terms(field string, values []any) Query {
	return Query{"terms": Query{field: values}}
}

func AggTerms(field string, size int) Query {
	return Query{
		"terms": Query{
			"field": field,
			"size":  size,
		},
	}
}

func AggAvg(field string) Query {
	return Query{
		"avg": Query{
			"field": field,
		},
	}
}

func AggSum(field string) Query {
	return Query{
		"sum": Query{
			"field": field,
		},
	}
}
