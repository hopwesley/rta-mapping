package common

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	"sync"
)

var _rtaMapping *RtaBitMap
var _rtaMapOnce sync.Once

type UserIDToRat map[int]*BitSlot

func RtaMapInst() *RtaBitMap {
	_rtaMapOnce.Do(func() {
		_rtaMapping = &RtaBitMap{rtaGrp: make(map[int64]bool)}
		for i := 0; i < BitSlotSize; i++ {
			_rtaMapping.bitMapGrp[i] = make(UserIDToRat)
		}
	})
	return _rtaMapping
}

type RtaBitMap struct {
	bitMapGrp  [BitSlotSize]UserIDToRat
	slotLock   [BitSlotSize]sync.RWMutex
	rtaGrp     map[int64]bool
	rtaTypLock sync.RWMutex
}

func (m *RtaBitMap) InitByOneRtaWithoutLock(rtaID int64, userIDs []int) error {
	var slotPos = rtaIDToSlotPos(rtaID)

	for _, usrID := range userIDs {
		var rtaMap = m.bitMapGrp[slotPos]
		if rtaMap[usrID] == nil {
			rtaMap[usrID] = new(BitSlot)
		}
		rtaMap[usrID].Add(slotPos)
	}

	m.rtaGrp[rtaID] = true

	return nil
}

type RtaUpdateItem struct {
	RtaID   int64 `json:"rta_id"`
	UserIDs []int `json:"user_ids"`
	IsDel   bool  `json:"is_del"`
}

func (m *RtaBitMap) UpdateRta(req *RtaUpdateItem) *UpdateResponse {
	slotPos := rtaIDToSlotPos(req.RtaID)

	for _, userID := range req.UserIDs {
		m.slotLock[slotPos].Lock()

		rtaMap := m.bitMapGrp[slotPos]
		if rtaMap[userID] == nil {
			rtaMap[userID] = new(BitSlot)
		}

		if req.IsDel {
			rtaMap[userID].Clear(slotPos)
		} else {
			rtaMap[userID].Add(slotPos)
		}

		m.slotLock[slotPos].Unlock()
	}

	m.rtaTypLock.Lock()
	m.rtaGrp[req.RtaID] = true
	m.rtaTypLock.Unlock()

	return &UpdateResponse{
		Success: true,
		Code:    0,
		Msg:     "Success",
	}
}

func (m *RtaBitMap) HintedList(userID int, ratIDs []int64) (result []int64) {

	for _, rtaID := range ratIDs {

		m.rtaTypLock.RLock()
		if m.rtaGrp[rtaID] == false {
			m.rtaTypLock.RUnlock()
			continue
		}
		m.rtaTypLock.RUnlock()

		slotPos := rtaIDToSlotPos(rtaID)
		m.slotLock[slotPos].RLock()

		rtaMap := m.bitMapGrp[slotPos]

		if rtaMap[userID] == nil {
			m.slotLock[slotPos].RUnlock()
			continue
		}
		if rtaMap[userID].Has(slotPos) {
			result = append(result, rtaID)
		}
		m.slotLock[slotPos].RUnlock()
	}
	return
}

func (m *RtaBitMap) HintedUserIDs(userID int, ratIDs []int64) (result []*UserInfo) {

	for _, rtaID := range ratIDs {

		m.rtaTypLock.RLock()
		if m.rtaGrp[rtaID] == false {
			m.rtaTypLock.RUnlock()
			continue
		}
		m.rtaTypLock.RUnlock()

		slotPos := rtaIDToSlotPos(rtaID)
		m.slotLock[slotPos].RLock()

		rtaMap := m.bitMapGrp[slotPos]
		if rtaMap[userID] == nil {
			m.slotLock[slotPos].RUnlock()
			continue
		}

		if rtaMap[userID].Has(slotPos) {
			ui := &UserInfo{
				RtaId:        rtaID,
				IsInterested: true,
			}
			result = append(result, ui)
		}
		m.slotLock[slotPos].RUnlock()
	}
	return
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
