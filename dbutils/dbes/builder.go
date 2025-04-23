package dbes

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Builder ES DSL构造器
type Builder struct {
	client *elasticsearch.Client
	body   map[string]interface{}
	index  string
	id     string
}

// NewBuilder 创建新的DSL构造器
func NewBuilder(client *elasticsearch.Client) *Builder {
	return &Builder{
		client: client,
		body:   make(map[string]interface{}),
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
	if m, ok := doc.(map[string]interface{}); ok {
		b.body = m
	} else {
		bytes, _ := json.Marshal(doc)
		json.Unmarshal(bytes, &b.body)
	}
	return b
}

// Build 构建DSL
func (b *Builder) Build() (map[string]interface{}, error) {
	return b.body, nil
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

// Search 执行搜索请求
func (b *Builder) Search(ctx context.Context) (*esapi.Response, error) {
	search := b.client.Search.WithContext(ctx)
	if b.index != "" {
		search = b.client.Search.WithIndex(b.index)
	}
	return b.client.Search(
		search,
		b.client.Search.WithBody(b.BuildReader()),
	)
}

// Count 执行count请求
func (b *Builder) Count(ctx context.Context) (*esapi.Response, error) {
	count := b.client.Count.WithContext(ctx)
	if b.index != "" {
		count = b.client.Count.WithIndex(b.index)
	}
	return b.client.Count(
		count,
		b.client.Count.WithBody(b.BuildReader()),
	)
}

// Create 创建文档
func (b *Builder) Create(ctx context.Context) (*esapi.Response, error) {
	if b.id != "" {
		req := esapi.CreateRequest{
			Index:      b.index,
			DocumentID: b.id,
			Body:       b.BuildReader(),
		}
		return req.Do(ctx, b.client)
	}
	req := esapi.IndexRequest{
		Index: b.index,
		Body:  b.BuildReader(),
	}
	return req.Do(ctx, b.client)
}

// Update 更新文档
func (b *Builder) Update(ctx context.Context) (*esapi.Response, error) {
	if b.id == "" {
		return nil, ErrIDRequired
	}
	req := esapi.UpdateRequest{
		Index:      b.index,
		DocumentID: b.id,
		Body:       b.BuildReader(),
	}
	return req.Do(ctx, b.client)
}

// UpdateByQuery 根据查询条件批量更新文档
func (b *Builder) UpdateByQuery(ctx context.Context) (*esapi.Response, error) {
	req := esapi.UpdateByQueryRequest{
		Index: []string{b.index},
		Body:  b.BuildReader(),
	}
	return req.Do(ctx, b.client)
}

// Delete 删除文档
func (b *Builder) Delete(ctx context.Context) (*esapi.Response, error) {
	if b.id == "" {
		return nil, ErrIDRequired
	}
	req := esapi.DeleteRequest{
		Index:      b.index,
		DocumentID: b.id,
	}
	return req.Do(ctx, b.client)
}

// DeleteByQuery 根据查询条件批量删除文档
func (b *Builder) DeleteByQuery(ctx context.Context) (*esapi.Response, error) {
	req := esapi.DeleteByQueryRequest{
		Index: []string{b.index},
		Body:  b.BuildReader(),
	}
	return req.Do(ctx, b.client)
}

// Bulk 批量操作
func (b *Builder) Bulk(ctx context.Context, actions []map[string]interface{}) (*esapi.Response, error) {
	var buf strings.Builder
	for _, action := range actions {
		actionLine, _ := json.Marshal(action)
		buf.Write(actionLine)
		buf.WriteString("\n")
	}
	req := esapi.BulkRequest{
		Index: b.index,
		Body:  strings.NewReader(buf.String()),
	}
	return req.Do(ctx, b.client)
}

// BulkBuilder 批量操作构造器
type BulkBuilder struct {
	actions []map[string]interface{}
}

// NewBulkBuilder 创建批量操作构造器
func NewBulkBuilder() *BulkBuilder {
	return &BulkBuilder{
		actions: make([]map[string]interface{}, 0),
	}
}

// Add 添加一个操作
func (bb *BulkBuilder) Add(action string, meta, doc map[string]interface{}) *BulkBuilder {
	bb.actions = append(bb.actions, map[string]interface{}{
		action: meta,
	})
	if doc != nil {
		bb.actions = append(bb.actions, doc)
	}
	return bb
}

// AddCreate 添加创建操作
func (bb *BulkBuilder) AddCreate(index, id string, doc interface{}) *BulkBuilder {
	return bb.Add("create", M("_index", index, "_id", id), toMap(doc))
}

// AddIndex 添加索引操作
func (bb *BulkBuilder) AddIndex(index, id string, doc interface{}) *BulkBuilder {
	return bb.Add("index", M("_index", index, "_id", id), toMap(doc))
}

// AddUpdate 添加更新操作
func (bb *BulkBuilder) AddUpdate(index, id string, doc interface{}) *BulkBuilder {
	return bb.Add("update", M("_index", index, "_id", id), M("doc", doc))
}

// AddDelete 添加删除操作
func (bb *BulkBuilder) AddDelete(index, id string) *BulkBuilder {
	return bb.Add("delete", M("_index", index, "_id", id), nil)
}

// Build 构建批量操作
func (bb *BulkBuilder) Build() []map[string]interface{} {
	return bb.actions
}

// M 用于构造map[string]interface{}的辅助方法
func M(pairs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(pairs); i += 2 {
		if i+1 < len(pairs) {
			m[pairs[i].(string)] = pairs[i+1]
		}
	}
	return m
}

// toMap 将结构体转换为map
func toMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	bytes, _ := json.Marshal(v)
	var result map[string]interface{}
	json.Unmarshal(bytes, &result)
	return result
}

// 一些常用的查询条件构造辅助方法
type Query struct{}

var Q Query

// Term 构造term查询
func (q Query) Term(field string, value interface{}) map[string]interface{} {
	return M("term", M(field, M("value", value)))
}

// Terms 构造terms查询
func (q Query) Terms(field string, values ...interface{}) map[string]interface{} {
	return M("terms", M(field, values))
}

// Match 构造match查询
func (q Query) Match(field string, value interface{}) map[string]interface{} {
	return M("match", M(field, value))
}

// Range 构造range查询
func (q Query) Range(field string, ranges map[string]interface{}) map[string]interface{} {
	return M("range", M(field, ranges))
}

// Bool 构造bool查询
func (q Query) Bool(must, should, mustNot, filter []map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
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
func (q Query) Script(script string, lang string, params map[string]interface{}) map[string]interface{} {
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
func (a Aggs) Terms(field string, size int) map[string]interface{} {
	return M("terms", M("field", field, "size", size))
}

// Avg 构造avg聚合
func (a Aggs) Avg(field string) map[string]interface{} {
	return M("avg", M("field", field))
}

// DateHistogram 构造date_histogram聚合
func (a Aggs) DateHistogram(field, interval string) map[string]interface{} {
	return M("date_histogram", M("field", field, "interval", interval))
}

// 一些常用的排序构造辅助方法
type Sort struct{}

var S Sort

// Field 构造字段排序
func (s Sort) Field(field, order string) map[string]interface{} {
	return M(field, M("order", order))
}

// Score 构造评分排序
func (s Sort) Score(order string) map[string]interface{} {
	return M("_score", M("order", order))
}

var (
	ErrIDRequired = errors.New("document ID is required")
)
