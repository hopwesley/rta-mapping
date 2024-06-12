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
	} else {
		m.Delete(key)
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

func (im *IDMap) UpdateIDMap(req *JsonRequest) *JsonResponse {
	im.updateMap(im.IMEIMD5, req.IMEIMD5, req.UserID)
	im.updateMap(im.OAID, req.OAID, req.UserID)
	im.updateMap(im.IDFA, req.IDFA, req.UserID)
	im.updateMap(im.AndroidIDMD5, req.AndroidIDMD5, req.UserID)
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
	im.IMEIMD5.Store(item.IMEIMD5, item.UserID)
	im.OAID.Store(item.OAID, item.UserID)
	im.IDFA.Store(item.IDFA, item.UserID)
	im.AndroidIDMD5.Store(item.AndroidIDMD5, item.UserID)
}
