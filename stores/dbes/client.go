package dbes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/morehao/go-tools/glog"
)

type ESConfig struct {
	Service  string `yaml:"service"`  // 服务名称
	Addr     string `yaml:"addr"`     // 地址
	User     string `yaml:"user"`     // 用户名
	Password string `yaml:"password"` // 密码
}

func InitES(cfg ESConfig) (*elasticsearch.Client, *elasticsearch.TypedClient, error) {
	customLogger, getLoggerErr := newEsLogger(&cfg)
	if getLoggerErr != nil {
		return nil, nil, getLoggerErr
	}
	commonCfg := elasticsearch.Config{
		Addresses: []string{cfg.Addr},
		Username:  cfg.User,
		Password:  cfg.Password,
		Logger:    customLogger,
	}
	simpleClient, newSimpleClientErr := elasticsearch.NewClient(commonCfg)
	if newSimpleClientErr != nil {
		return nil, nil, newSimpleClientErr
	}
	typedClient, newTypedClientErr := elasticsearch.NewTypedClient(commonCfg)
	if newTypedClientErr != nil {
		return nil, nil, newTypedClientErr
	}
	return simpleClient, typedClient, nil
}

func newEsLogger(cfg *ESConfig) (*esLog, error) {
	l, err := glog.GetLogger()
	if err != nil {
		return nil, err
	}
	return &esLog{
		logger:  l,
		service: cfg.Service,
	}, nil
}

type esLog struct {
	logger  glog.Logger
	service string
}

func (l *esLog) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	ctx := req.Context()
	cost := dur.Nanoseconds() / 1e4 / 100.0
	end := start.Add(dur)

	// 获取查询的 HTTP method 和路径
	method := req.Method
	path := req.URL.Path
	realCode := res.StatusCode

	var fields []any
	fields = append(fields,
		glog.KeyService, l.service,
		glog.KeyProto, glog.ValueProtoES,
		glog.KeyRequestStartTime, glog.FormatRequestTime(start),
		glog.KeyRequestEndTime, glog.FormatRequestTime(end),
		glog.KeyCost, cost,
		glog.KeyRalCode, realCode,
		glog.KeyDslMethod, method,
		glog.KeyDslPath, path,
	)
	msg := "es execute success"
	if err != nil {
		realCode = -1
		msg = err.Error()
		fields = append(fields, glog.KeyErrorMsg, msg)
		l.logger.Errorw(ctx, msg, fields...)
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
	var affectedRows int
	if res.Body != nil && res.Body != http.NoBody {
		bodyBytes, readErr := io.ReadAll(res.Body)
		defer func() {
			res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}()
		if readErr != nil {
			l.logger.Errorw(ctx, fmt.Sprintf("read es response body fail, error: %s", readErr.Error()), fields...)
			return nil
		}

		var resBody map[string]any
		if err := jsoniter.Unmarshal(bodyBytes, &resBody); err == nil {
			if hits, ok := resBody["hits"].(map[string]any); ok {
				if total, ok := hits["total"].(map[string]any); ok {
					if value, ok := total["value"].(float64); ok {
						affectedRows = int(value)
					}
				}
			}
		}
		fields = append(fields,
			glog.KeyAffectedRows, affectedRows,
		)
	}
	if realCode != 200 {
		l.logger.Errorw(ctx, msg, fields...)
	} else {
		l.logger.Infow(ctx, msg, fields...)
	}
	return err
}

func (l *esLog) RequestBodyEnabled() bool {
	return true
}

func (l *esLog) ResponseBodyEnabled() bool {
	return true
}
