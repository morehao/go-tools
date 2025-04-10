package dbClient

import (
	"bytes"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/morehao/go-tools/glog"
	"go.uber.org/zap"
)

type ESConfig struct {
	Service  string `yaml:"service"`  // 服务名称
	Addr     string `yaml:"addr"`     // 地址
	User     string `yaml:"user"`     // 用户名
	Password string `yaml:"password"` // 密码
}

func InitES(cfg ESConfig) (*elasticsearch.TypedClient, error) {
	customLogger, getLoggerErr := newEsLogger(&cfg)
	if getLoggerErr != nil {
		return nil, getLoggerErr
	}
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{cfg.Addr},
		Username:  cfg.User,
		Password:  cfg.Password,
		Logger:    customLogger,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newEsLogger(cfg *ESConfig) (*esLog, error) {
	l, err := glog.GetLogger(glog.WithZapOptions(zap.AddCallerSkip(2)))
	if err != nil {
		return nil, err
	}
	return &esLog{
		logger: l,
	}, nil
}

type esLog struct {
	logger glog.Logger
}

func (l *esLog) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	ctx := req.Context()
	cost := dur.Nanoseconds() / 1e4 / 100.0
	end := start.Add(dur)

	// 假设通过解析响应体获取生效行数 (以 Elasticsearch 为例)
	var affectedRows int
	if res.Body != nil {
		var resBody map[string]interface{}
		decoder := jsoniter.NewDecoder(res.Body)
		if err := decoder.Decode(&resBody); err == nil {
			// 假设查询结果中有 hits.total.value 字段
			if hits, ok := resBody["hits"].(map[string]interface{}); ok {
				if total, ok := hits["total"].(map[string]interface{}); ok {
					if value, ok := total["value"].(float64); ok {
						affectedRows = int(value)
					}
				}
			}
		}
	}

	// 获取查询的 HTTP method 和路径
	method := req.Method
	path := req.URL.Path
	realCode := res.StatusCode

	var fields []any
	fields = append(fields,
		glog.KeyRequestStartTime, glog.FormatRequestTime(start),
		glog.KeyRequestEndTime, glog.FormatRequestTime(end),
		glog.KeyCost, cost,
		glog.KeyRalCode, realCode, // 添加状态码
		glog.KeyAffectedRows, affectedRows, // 添加生效行数
		"dslMethod", method, // 添加请求方法
		"dslPath", path, // 添加请求路径
	)
	msg := "es execute success"
	if err != nil {
		realCode = -1
		msg = err.Error()
		fields = append(fields, glog.KeyErrorMsg, msg)
	}

	if req.Body != nil && req.Body != http.NoBody {
		var buf bytes.Buffer
		if req.GetBody != nil {
			b, _ := req.GetBody()
			buf.ReadFrom(b)
		} else {
			buf.ReadFrom(req.Body)
		}
		fields = append(fields, glog.KeyDsl, buf.String())
	}
	if realCode != 200 {
		l.logger.Errorw(ctx, msg, fields...)
	} else {
		l.logger.Infow(ctx, msg, fields...)
	}
	return err
}

func (l *esLog) RequestBodyEnabled() bool {
	return false
}

func (l *esLog) ResponseBodyEnabled() bool {
	return false
}
