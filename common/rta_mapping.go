package common

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sync"
)

var _rtaMapping *RtaBitMap
var _rtaMapOnce sync.Once

func RtaMapInst() *RtaBitMap {
	_rtaMapOnce.Do(func() {
		_rtaMapping = &RtaBitMap{
			rtaGrp:    new(sync.Map),
			bitMapGrp: new(sync.Map),
		}
	})

	return _rtaMapping
}

type RtaBitMap struct {
	bitMapGrp *sync.Map
	rtaGrp    *sync.Map
}

func (m *RtaBitMap) InitByOneRtaWithoutLock(rtaID int64, userIDs []int) {
	var slotPos = rtaIDToSlotPos(rtaID)
	for _, userID := range userIDs {
		bitSlot, _ := m.bitMapGrp.LoadOrStore(userID, new(BitSlot))
		bitSlot.(*BitSlot).Add(slotPos)
	}
	m.rtaGrp.Store(slotPos, rtaID)
}

type RtaUpdateItem struct {
	RtaID   int64 `json:"rta_id"`
	UserIDs []int `json:"user_ids"`
	IsDel   bool  `json:"is_del"`
}

func (m *RtaBitMap) UpdateRta(req *RtaUpdateItem) *JsonResponse {

	slotPos := rtaIDToSlotPos(req.RtaID)

	for _, userID := range req.UserIDs {
		bitSlot, _ := m.bitMapGrp.LoadOrStore(userID, new(BitSlot))
		if req.IsDel {
			bitSlot.(*BitSlot).Clear(slotPos)
		} else {
			bitSlot.(*BitSlot).Add(slotPos)
		}
	}
	m.rtaGrp.Store(slotPos, req.RtaID)
	return SuccessJsonRes
}

func (m *RtaBitMap) HintedList(userID int, ratIDs []int64) (result []int64) {

	value, ok := m.bitMapGrp.Load(userID)
	if !ok {
		return
	}
	bitMap := value.(*BitSlot)

	for _, rtaID := range ratIDs {
		slotPos := rtaIDToSlotPos(rtaID)

		value, ok := m.rtaGrp.Load(slotPos)
		if !ok || value.(int64) <= 0 {
			continue
		}

		if bitMap.Has(slotPos) {
			result = append(result, rtaID)
		}
	}
	return
}

func (m *RtaBitMap) HintedUserIDs(userID int, ratIDs []int64) (result []*UserInfo) {

	value, ok := m.bitMapGrp.Load(userID)
	if !ok {
		return
	}
	bitMap := value.(*BitSlot)

	for _, rtaID := range ratIDs {
		slotPos := rtaIDToSlotPos(rtaID)

		value, ok := m.rtaGrp.Load(slotPos)
		if !ok || value.(int64) <= 0 {
			continue
		}

		if bitMap.Has(slotPos) {
			ui := &UserInfo{
				RtaId:        rtaID,
				IsInterested: true,
			}
			result = append(result, ui)
		}
	}

	return
}

func (m *RtaBitMap) QueryRatInfos(request *JsonRequest) *JsonResponse {
	value, ok := m.bitMapGrp.Load(request.UserID)
	if !ok {
		return NotFoundJsonRes
	}
	data := value.(*BitSlot)

	var posArr = data.GetBitsAsArray()
	var rtaIDs []int64
	for _, pos := range posArr {
		if rtaID, ok := m.rtaGrp.Load(pos); ok {
			rtaIDs = append(rtaIDs, rtaID.(int64))
		}
	}

	bts, _ := json.Marshal(rtaIDs)
	return &JsonResponse{
		Success: true,
		Msg:     string(bts),
	}
}

func CheckIfHinted(request *Req) *Rsp {
	userID, ok := IDMapInst().DeviceToUserID(request.Device)
	if !ok {
		return FailureHit(request.ReqId)
	}
	var rtaIDs = request.RtaIds
	if len(rtaIDs) < 1 {
		return FailureHit(request.ReqId)
	}

	var uis = RtaMapInst().HintedUserIDs(userID, rtaIDs)
	if len(uis) < 1 {
		return FailureHit(request.ReqId)
	}

	var response = &Rsp{
		StatusCode: HitSuccess,
		BidType:    &wrapperspb.Int32Value{Value: BidTypeOk},
		UserInfos:  uis,
		ReqId:      request.ReqId,
	}

	return response
}
