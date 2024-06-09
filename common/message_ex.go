package common

var (
	HitFailure = Response{
		StatusCode: HitFailed,
	}
)

func FailureHit(reqID string) *Response {
	return &Response{
		StatusCode: HitFailed,
		ReqId:      reqID,
	}
}
