package restful

type ErrCode int

const (
	Success ErrCode = iota
	RequestError
	SearchError
)

type Restful[T any] struct {
	Code ErrCode `json:"code"`
	Msg  string  `json:"msg"`
	Data T       `json:"data"`
}
