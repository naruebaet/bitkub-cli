package templ

const Response = `
package model

import (
	"time"
)

type BaseResponse struct {
	Code int         ` + "`json:\"code\"`" + `
	Data interface{} ` + "`json:\"data\"`" + `
}

type BaseErrorResponse struct {
	Code     int    ` + "`json:\"code\"`" + `
	ErrTitle string ` + "`json:\"error\"`" + `
}

func NewBaseResponse(code int, data interface{}) *BaseResponse {
	if data == nil {
		type Empty struct{}
		data = Empty{}
	}
	return &BaseResponse{
		Code: code,
		Data: data,
	}
}

func NewBaseErrorResponse(code int, title string) *BaseErrorResponse {
	return &BaseErrorResponse{
		Code:     code,
		ErrTitle: title,
	}
}

// VerifyRecaptchaResponse struct
type VerifyRecaptchaResponse struct {
	Success     bool      ` + "`json:\"success\"`" + `
	Score       float64   ` + "`json:\"score\"`" + `
	Action      string   ` + "`json:\"action\"`" + `
	ChallengeTS time.Time` + "`json:\"challenge_ts\"`" + `
	Hostname    string   ` + "`json:\"hostname\"`" + `
	ErrorCodes  []string ` + "`json:\"error-codes\"`" + `
}
`
