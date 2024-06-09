package workV1

import (
	"github.com/hopwesley/rta-mapping/common"
)

func QueryRtaMap(reqBody *common.Request) *common.Response {

	var response = &common.Response{
		StatusCode: common.HitSuccess,
		BidType:    common.BidTypeOk,
		UserInfos:  nil,
	}

	return response
}
