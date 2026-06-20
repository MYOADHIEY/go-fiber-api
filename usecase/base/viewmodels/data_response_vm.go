package basevm

type BaseResponseErrorVM struct {
	Messages interface{} `json:"messages"`
}

type BaseResponseVM struct {
	StatusCode int         `json:"code"`
	Data       interface{} `json:"data"`
	Meta       interface{} `json:"meta"`
	Errors     []string    `json:"errors"`
}
