package common

var (
	HitFailure = Rsp{
		StatusCode: HitFailed,
	}
)

func FailureHit(reqID string) *Rsp {
	return &Rsp{
		StatusCode: HitFailed,
		ReqId:      reqID,
	}
}
