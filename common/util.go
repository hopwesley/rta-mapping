package common

var (
	Version   string
	Commit    string
	BuildTime string
)

const (
	HitSuccess = iota
	HitFailed
)
const (
	BidTypeOk = iota
)

var (
	HitFailure = &Response{
		StatusCode: HitFailed,
	}
)
