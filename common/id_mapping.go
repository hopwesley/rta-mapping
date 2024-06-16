package common

import (
	"strconv"
	"sync"
)

// 所以把imei放最前面，其次oaid，然后idfa

type IDMap struct {
	IMEIMD5      *sync.Map
	OAID         *sync.Map
	IDFA         *sync.Map
	AndroidIDMD5 *sync.Map
}

const (
	IDOpTypAdd = iota
	IDOpTypUpdate
	IDOpTypDel
)

type IDOpItem struct {
	Val   string `json:"val"`
	OpTyp int    `json:"op_typ"`
}
type IDUpdateReq struct {
	UserID       int       `json:"user_id,omitempty"`
	IMEIMD5      *IDOpItem `json:"imei_md5,omitempty"`
	OAID         *IDOpItem `json:"oaid,omitempty"`
	IDFA         *IDOpItem `json:"idfa,omitempty"`
	AndroidIDMD5 *IDOpItem `json:"android_id_md5,omitempty"`
}

var _idMapping *IDMap
var _idMapOnce sync.Once

func IDMapInst() *IDMap {
	_idMapOnce.Do(func() {
		_idMapping = &IDMap{
			IMEIMD5:      new(sync.Map),
			OAID:         new(sync.Map),
			IDFA:         new(sync.Map),
			AndroidIDMD5: new(sync.Map),
		}
	})
	return _idMapping
}

func (im *IDMap) getFromMap(m *sync.Map, key string) (int, bool) {
	value, ok := m.Load(key)
	if !ok {
		return -1, false
	}
	return value.(int), true
}

func (im *IDMap) updateMap(m *sync.Map, key string, userID int) {
	if len(key) > 0 {
		m.Store(key, userID)
	}
}

func (im *IDMap) DeviceToUserID(device *Device) (int, bool) {

	if userID, ok := im.getFromMap(im.IMEIMD5, device.ImeiMd5); ok {
		return userID, true
	}

	if userID, ok := im.getFromMap(im.OAID, device.Oaid); ok {
		return userID, true
	}

	if userID, ok := im.getFromMap(im.IDFA, device.Idfa); ok {
		return userID, true
	}

	if userID, ok := im.getFromMap(im.AndroidIDMD5, device.AndroidIdMd5); ok {
		return userID, true
	}

	return -1, false
}
func operateID(m *sync.Map, item *IDOpItem, userID int) {
	if item == nil {
		return
	}

	switch item.OpTyp {
	case IDOpTypAdd, IDOpTypUpdate:
		m.Store(item.Val, userID)
		break
	case IDOpTypDel:
		m.Delete(item.Val)
		break
	}
}

func (im *IDMap) UpdateIDMap(req []*IDUpdateReq) *JsonResponse {
	for _, idOp := range req {
		operateID(im.IMEIMD5, idOp.IMEIMD5, idOp.UserID)
		operateID(im.OAID, idOp.OAID, idOp.UserID)
		operateID(im.IDFA, idOp.IDFA, idOp.UserID)
		operateID(im.AndroidIDMD5, idOp.AndroidIDMD5, idOp.UserID)
	}
	return SuccessJsonRes
}

func (im *IDMap) QueryIDInfos(req *JsonRequest) *JsonResponse {

	var result = &JsonResponse{
		Success: true,
	}
	if userID, ok := im.getFromMap(im.IMEIMD5, req.IMEIMD5); ok {
		result.Msg = strconv.Itoa(userID)
		return result
	}

	if userID, ok := im.getFromMap(im.OAID, req.OAID); ok {
		result.Msg = strconv.Itoa(userID)
		return result
	}

	if userID, ok := im.getFromMap(im.IDFA, req.IDFA); ok {
		result.Msg = strconv.Itoa(userID)
		return result
	}

	if userID, ok := im.getFromMap(im.AndroidIDMD5, req.AndroidIDMD5); ok {
		result.Msg = strconv.Itoa(userID)
		return result
	}

	return NotFoundJsonRes
}

func (im *IDMap) UpdateByMySqlWithoutLock(item JsonRequest) {
	im.updateMap(im.IMEIMD5, item.IMEIMD5, item.UserID)
	im.updateMap(im.OAID, item.OAID, item.UserID)
	im.updateMap(im.IDFA, item.IDFA, item.UserID)
	//fmt.Println("\n\n", item.IDFA, item.UserID)
	//fmt.Println(im.getFromMap(im.IDFA, item.IDFA))
	//fmt.Println()
	im.updateMap(im.AndroidIDMD5, item.AndroidIDMD5, item.UserID)
}
