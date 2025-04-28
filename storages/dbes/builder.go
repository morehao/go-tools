package dbes

import (
	"encoding/json"
	"strings"
)

// Builder ES DSL构造器
type Builder struct {
	body  map[string]any
	index string
	id    string
}

// NewBuilder 创建新的DSL构造器
func NewBuilder() *Builder {
	return &Builder{
		body: make(map[string]any),
	}
}

// Set 通用设置方法，支持链式调用
func (b *Builder) Set(key string, value interface{}) *Builder {
	b.body[key] = value
	return b
}

// SetQuery 设置查询条件
func (b *Builder) SetQuery(query interface{}) *Builder {
	return b.Set("query", query)
}

// SetAggs 设置聚合
func (b *Builder) SetAggs(aggs interface{}) *Builder {
	return b.Set("aggs", aggs)
}

// SetSort 设置排序
func (b *Builder) SetSort(sort interface{}) *Builder {
	return b.Set("sort", sort)
}

// SetSize 设置返回数量
func (b *Builder) SetSize(size int) *Builder {
	return b.Set("size", size)
}

// SetFrom 设置偏移量
func (b *Builder) SetFrom(from int) *Builder {
	return b.Set("from", from)
}

// SetSource 设置返回字段
func (b *Builder) SetSource(source interface{}) *Builder {
	return b.Set("_source", source)
}

// SetIndex 设置要操作的索引
func (b *Builder) SetIndex(index string) *Builder {
	b.index = index
	return b
}

// SetID 设置文档ID
func (b *Builder) SetID(id string) *Builder {
	b.id = id
	return b
}

// SetDoc 设置文档内容
func (b *Builder) SetDoc(doc interface{}) *Builder {
	if m, ok := doc.(map[string]any); ok {
		b.body = m
	} else {
		bytes, _ := json.Marshal(doc)
		json.Unmarshal(bytes, &b.body)
	}
	return b
}

// Build 构建DSL
func (b *Builder) Build() map[string]any {
	return b.body
}

// BuildBytes 构建并返回JSON字节数组
func (b *Builder) BuildBytes() ([]byte, error) {
	return json.Marshal(b.body)
}

// BuildReader 构建io.Reader用于ES请求
func (b *Builder) BuildReader() *strings.Reader {
	body, _ := b.BuildBytes()
	return strings.NewReader(string(body))
}

// GetIndex 获取索引
func (b *Builder) GetIndex() string {
	return b.index
}

// GetID 获取文档ID
func (b *Builder) GetID() string {
	return b.id
}

// M 用于构造map[string]interface{}的辅助方法
func M(pairs ...interface{}) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			m[pairs[i].(string)] = pairs[i+1]
		}
	}
	return m
}

// toMap 将结构体转换为map
func toMap(v interface{}) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	bytes, _ := json.Marshal(v)
	var result map[string]any
	json.Unmarshal(bytes, &result)
	return result
}

// 一些常用的查询条件构造辅助方法
type Query struct{}

var Q Query

// Term 构造term查询
func (q Query) Term(field string, value interface{}) map[string]any {
	return M("term", M(field, M("value", value)))
}

// Terms 构造terms查询
func (q Query) Terms(field string, values ...interface{}) map[string]any {
	return M("terms", M(field, values))
}

// Match 构造match查询
func (q Query) Match(field string, value interface{}) map[string]any {
	return M("match", M(field, value))
}

// Range 构造range查询
func (q Query) Range(field string, ranges map[string]any) map[string]any {
	return M("range", M(field, ranges))
}

// Bool 构造bool查询
func (q Query) Bool(must, should, mustNot, filter []map[string]any) map[string]any {
	m := make(map[string]any)
	if len(must) > 0 {
		m["must"] = must
	}
	if len(should) > 0 {
		m["should"] = should
	}
	if len(mustNot) > 0 {
		m["must_not"] = mustNot
	}
	if len(filter) > 0 {
		m["filter"] = filter
	}
	return M("bool", m)
}

// Script 构造脚本查询
func (q Query) Script(script string, lang string, params map[string]any) map[string]any {
	scriptMap := M("source", script)
	if lang != "" {
		scriptMap["lang"] = lang
	}
	if params != nil {
		scriptMap["params"] = params
	}
	return M("script", scriptMap)
}

// 一些常用的聚合构造辅助方法
type Aggs struct{}

var A Aggs

// Terms 构造terms聚合
func (a Aggs) Terms(field string, size int) map[string]any {
	return M("terms", M("field", field, "size", size))
}

// Avg 构造avg聚合
func (a Aggs) Avg(field string) map[string]any {
	return M("avg", M("field", field))
}

// DateHistogram 构造date_histogram聚合
func (a Aggs) DateHistogram(field, interval string) map[string]any {
	return M("date_histogram", M("field", field, "interval", interval))
}

// 一些常用的排序构造辅助方法
type Sort struct{}

var S Sort

// Field 构造字段排序
func (s Sort) Field(field, order string) map[string]any {
	return M(field, M("order", order))
}

// Score 构造评分排序
func (s Sort) Score(order string) map[string]any {
	return M("_score", M("order", order))
}
