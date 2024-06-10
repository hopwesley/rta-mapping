package common

import (
	"sync"
)

// 所以把imei放最前面，其次oaid，然后idfa

type IDMap struct {
	IMEIMD5      map[string]string `json:"imei_md5"`
	OAID         map[string]string `json:"oaid"`
	IDFA         map[string]string `json:"idfa"`
	AndroidIDMD5 map[string]string `json:"android_id_md5"`
}

var _idMapping *IDMap
var _idMapOnce sync.Once

func IDMapInst() *IDMap {
	_idMapOnce.Do(func() {
		_idMapping = &IDMap{
			IMEIMD5:      make(map[string]string),
			OAID:         make(map[string]string),
			IDFA:         make(map[string]string),
			AndroidIDMD5: make(map[string]string),
		}
	})
	return _idMapping
}

func (im *IDMap) QueryKey(device *Device) string {
	phoneMD5, ok := im.IMEIMD5[device.ImeiMd5]
	if ok {
		return phoneMD5
	}
	phoneMD5, ok = im.OAID[device.Oaid]
	if ok {
		return phoneMD5
	}
	phoneMD5, ok = im.IDFA[device.Idfa]
	if ok {
		return phoneMD5
	}
	phoneMD5, ok = im.AndroidIDMD5[device.AndroidIdMd5]
	if ok {
		return phoneMD5
	}
	return ""
}

type IDUpdateRequest struct {
	UserID       string `json:"phone"`
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
