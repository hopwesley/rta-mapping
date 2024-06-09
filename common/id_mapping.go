package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"sync"
)

// 所以把imei放最前面，其次oaid，然后idfa

type IDMap struct {
	IMEIMD5      map[string]string `json:"imei_md5"`
	OAID         map[string]string `json:"oaid"`
	IDFA         map[string]string `json:"idfa"`
	AndroidIDMD5 map[string]string `json:"android_id_md5"`
}

var _mapping *IDMap
var _mapOnce sync.Once

func IDMapInst() *IDMap {
	_mapOnce.Do(func() {
		_mapping = &IDMap{
			IMEIMD5:      make(map[string]string),
			OAID:         make(map[string]string),
			IDFA:         make(map[string]string),
			AndroidIDMD5: make(map[string]string),
		}
	})
	return _mapping
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
	Phone        string `json:"phone"`
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
	h := md5.New()
	_, _ = io.WriteString(h, req.Phone)
	md5sum := h.Sum(nil)

	phoneMD5 := fmt.Sprintf("%x", md5sum)

	if len(req.IMEIMD5) > 0 {
		im.IMEIMD5[req.IMEIMD5] = phoneMD5
	}
	if len(req.OAID) > 0 {
		im.OAID[req.OAID] = phoneMD5
	}
	if len(req.IDFA) > 0 {
		im.IDFA[req.IDFA] = phoneMD5
	}
	if len(req.AndroidIDMD5) > 0 {
		im.AndroidIDMD5[req.AndroidIDMD5] = phoneMD5
	}

	return &UpdateResponse{
		Success: true,
		Code:    0,
		Msg:     "Success",
	}
}
