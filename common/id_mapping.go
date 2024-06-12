package common

import (
	"sync"
)

// 所以把imei放最前面，其次oaid，然后idfa

type IDMap struct {
	sync.RWMutex
	IMEIMD5      map[string]int `json:"imei_md5"`
	OAID         map[string]int `json:"oaid"`
	IDFA         map[string]int `json:"idfa"`
	AndroidIDMD5 map[string]int `json:"android_id_md5"`
}

var _idMapping *IDMap
var _idMapOnce sync.Once

func IDMapInst() *IDMap {
	_idMapOnce.Do(func() {
		_idMapping = &IDMap{
			IMEIMD5:      make(map[string]int),
			OAID:         make(map[string]int),
			IDFA:         make(map[string]int),
			AndroidIDMD5: make(map[string]int),
		}
	})
	return _idMapping
}

func (im *IDMap) DeviceToUserID(device *Device) (int, bool) {
	im.RLock()
	defer im.RUnlock()

	userID, ok := im.IMEIMD5[device.ImeiMd5]
	if ok {
		return userID, true
	}
	userID, ok = im.OAID[device.Oaid]
	if ok {
		return userID, true
	}
	userID, ok = im.IDFA[device.Idfa]
	if ok {
		return userID, true
	}
	userID, ok = im.AndroidIDMD5[device.AndroidIdMd5]
	if ok {
		return userID, true
	}
	return -1, false
}

type IDUpdateRequest struct {
	UserID       int    `json:"user_id"`
	IMEIMD5      string `json:"imei_md5"`
	OAID         string `json:"oaid"`
	IDFA         string `json:"idfa"`
	AndroidIDMD5 string `json:"android_id_md5"`
}

type UpdateResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
}

func (im *IDMap) UpdateIDMap(req *IDUpdateRequest) *UpdateResponse {
	im.Lock()
	defer im.Unlock()

	if len(req.IMEIMD5) > 0 {
		im.IMEIMD5[req.IMEIMD5] = req.UserID
	} else {
		delete(im.IMEIMD5, req.IMEIMD5)
	}
	if len(req.OAID) > 0 {
		im.OAID[req.OAID] = req.UserID
	} else {
		delete(im.OAID, req.OAID)
	}
	if len(req.IDFA) > 0 {
		im.IDFA[req.IDFA] = req.UserID
	} else {
		delete(im.IDFA, req.IDFA)
	}

	if len(req.AndroidIDMD5) > 0 {
		im.AndroidIDMD5[req.AndroidIDMD5] = req.UserID
	} else {
		delete(im.AndroidIDMD5, req.AndroidIDMD5)
	}

	return &UpdateResponse{
		Success: true,
		Code:    0,
		Msg:     "Success",
	}
}

func (im *IDMap) CleanMap() {
	im.IDFA = make(map[string]int)
	im.IMEIMD5 = make(map[string]int)
	im.OAID = make(map[string]int)
	im.AndroidIDMD5 = make(map[string]int)
}

func (im *IDMap) UpdateByMySqlWithoutLock(item IDUpdateRequest) {
	im.IDFA[item.IDFA] = item.UserID
	im.IMEIMD5[item.IMEIMD5] = item.UserID
	im.OAID[item.OAID] = item.UserID
	im.AndroidIDMD5[item.AndroidIDMD5] = item.UserID
}
