package dbes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/morehao/go-tools/glog"
)

func newEsLogger(cfg *ESConfig) (*esLog, error) {
	l, err := glog.GetLogger(cfg.loggerConfig, glog.WithCallerSkip(8))
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
	ralCode := res.StatusCode

	var fields []any
	fields = append(fields,
		glog.KeyService, l.service,
		glog.KeyProto, glog.ValueProtoES,
		glog.KeyRequestStartTime, glog.FormatRequestTime(start),
		glog.KeyRequestEndTime, glog.FormatRequestTime(end),
		glog.KeyCost, cost,
		glog.KeyRalCode, ralCode,
		glog.KeyDslMethod, method,
		glog.KeyDslPath, path,
	)
	msg := "es execute success"
	if err != nil {
		ralCode = -1
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

		if ralCode != 200 {
			msg = string(bodyBytes)
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
	if ralCode != 200 {
		l.logger.Errorw(ctx, msg, fields...)
	} else {
		l.logger.Debugw(ctx, msg, fields...)
	}
	return err
}

func (l *esLog) RequestBodyEnabled() bool {
	return true
}

func (l *esLog) ResponseBodyEnabled() bool {
	return true
}
