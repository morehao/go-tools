package esquery

import jsoniter "github.com/json-iterator/go"

type Query map[string]any

type BoolQuery struct {
	MustClauses    []Query
	ShouldClauses  []Query
	MustNotClauses []Query
	FilterClauses  []Query
	SortFields     []Query
	FromVal        *int
	SizeVal        *int
	AggsMap        Query
	Options        []func(Query) // 用于存储扩展选项
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

func (b *BoolQuery) Sort(field string, asc bool) *BoolQuery {
	order := "asc"
	if !asc {
		order = "desc"
	}
	b.SortFields = append(b.SortFields, Query{field: Query{"order": order}})
	return b
}

func (b *BoolQuery) From(from int) *BoolQuery {
	b.FromVal = &from
	return b
}

func (b *BoolQuery) Size(size int) *BoolQuery {
	b.SizeVal = &size
	return b
}

func (b *BoolQuery) Aggs(name string, q Query) *BoolQuery {
	if b.AggsMap == nil {
		b.AggsMap = Query{}
	}
	b.AggsMap[name] = q
	return b
}

// With 用于添加自定义扩展，用户可以通过传入函数进行扩展
func (b *BoolQuery) With(option func(Query)) *BoolQuery {
	b.Options = append(b.Options, option)
	return b
}

// 构建 bool 查询部分
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

// 构建完整的 DSL 请求体
func (b *BoolQuery) BuildDSL() Query {
	// 首先构建 query 部分
	query := b.BuildQuery()

	// 构建完整的查询体
	result := Query{
		"query": query, // 包含 query 部分
	}

	// 将排序、分页等移至查询体的根部
	if len(b.SortFields) > 0 {
		result["sort"] = b.SortFields
	}
	if b.FromVal != nil {
		result["from"] = *b.FromVal
	}
	if b.SizeVal != nil {
		result["size"] = *b.SizeVal
	}
	if len(b.AggsMap) > 0 {
		result["aggs"] = b.AggsMap
	}

	// 执行扩展选项
	for _, opt := range b.Options {
		opt(result)
	}

	// 返回完整的 DSL 查询
	return result
}

// 序列化为字符串
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

// 聚合构造器
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
