package workV1

import (
	"github.com/hopwesley/rta-mapping/common"
)

func QueryRtaMap(request *common.Request) *common.Response {

	var phoneMD5 = common.IDMapInst().QueryKey(request.Device)
	if len(phoneMD5) < 4 {
		return common.FailureHit(request.ReqId)
	}
	var rtaIDs = request.RtaIds
	var result = IDRatInst().HintList(phoneMD5, rtaIDs)
	if len(result) < 1 {
		return common.FailureHit(request.ReqId)
	}

	var uis []*common.UserInfos
	for _, ratId := range result {
		ui := &common.UserInfos{
			RtaId:        ratId,
			IsInterested: true,
		}
		uis = append(uis, ui)
	}

	var response = &common.Response{
		StatusCode: common.HitSuccess,
		BidType:    common.BidTypeOk,
		UserInfos:  uis,
		ReqId:      request.ReqId,
	}

	return response
}
