package V1

import (
	"fmt"
	"github.com/hopwesley/rta-mapping/common"
	"sync"
)

var _mapping *IDRatMap
var _mapOnce sync.Once

type IDRatMap struct {
	Relation map[string]bool
}

func IDRatInst() *IDRatMap {
	_mapOnce.Do(func() {
		_mapping = &IDRatMap{
			Relation: make(map[string]bool),
		}
	})

	return _mapping
}

func (ir *IDRatMap) Hint(key string) bool {
	return ir.Relation[key]
}

func (ir *IDRatMap) HintList(phoneMD5 string, ratIDs []int64) (result []int64) {

	for _, rid := range ratIDs {
		key := fmt.Sprintf("%s%d", phoneMD5, rid)
		if ir.Relation[key] {
			result = append(result, rid)
		}
	}

	return
}

type RtaUpdateItem struct {
	RtaID    int64  `json:"rta_id"`
	PhoneMD5 string `json:"phone_md5"`
	IsDel    bool   `json:"is_del"`
}

type RtaUpdateRequest struct {
	OPList []RtaUpdateItem `json:"op_list"`
}

func (ir *IDRatMap) UpdateHintList(req *RtaUpdateRequest) *common.UpdateResponse {

	for _, item := range req.OPList {
		key := fmt.Sprintf("%s%d", item.PhoneMD5, item.RtaID)
		if item.IsDel {
			delete(ir.Relation, key)
		} else {
			ir.Relation[key] = true
		}
	}

	return &common.UpdateResponse{
		Success: true,
		Code:    0,
		Msg:     "Success",
	}
}
