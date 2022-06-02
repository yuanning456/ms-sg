package common

type Result struct {
	Code int `json:"code"`
	Errmsg string `json:"errmsg"`
	Data interface{} `json:"data"`
}

func Error(code int,msg string) *Result  {
	return &Result{
		Code: code,
		Errmsg: msg,
	}
}

func Success(code int,data interface{}) *Result  {
	return &Result{
		Code: code,
		Data: data,
	}
}
