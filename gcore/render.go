package gcore

// ResponseRender 返回数据格式化
type ResponseRender interface {
	SetCode(int)
	SetMsg(string)
	SetData(interface{})
	SetDataWithFormat(interface{})
}

func NewResponseRender() ResponseRender {
	return &responseRender{}
}

type responseRender struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *responseRender) SetCode(code int) {
	r.Code = code
}
func (r *responseRender) SetMsg(msg string) {
	r.Msg = msg
}
func (r *responseRender) SetData(data interface{}) {
	r.Data = data
}

func (r *responseRender) SetDataWithFormat(data interface{}) {
	ResponseFormat(data)
	r.Data = data
}
