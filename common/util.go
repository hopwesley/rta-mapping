package common

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"hash/fnv"
	"io"
	"net/http"
	"runtime"
)

var (
	Version   string
	Commit    string
	BuildTime string
)
var (
	SuccessJsonRes = &JsonResponse{
		Success: true,
		Code:    0,
		Msg:     "Success",
	}
	NotFoundJsonRes = &JsonResponse{
		Success: false,
		Code:    -1,
		Msg:     "Not Found",
	}
)

const (
	HitSuccess = iota
	HitFailed
)
const (
	BidTypeOk = iota
)

type JsonResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
}

type JsonRequest struct {
	UserID       int    `json:"user_id"`
	IMEIMD5      string `json:"imei_md5"`
	OAID         string `json:"oaid"`
	IDFA         string `json:"idfa"`
	AndroidIDMD5 string `json:"android_id_md5"`
}

func ReadProtoRequest(w http.ResponseWriter, r *http.Request) *Req {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return nil
	}

	var request = &Req{}
	if err := proto.Unmarshal(body, request); err != nil {
		LogInst().Error("request is invalid:", err)
		return nil
	}

	return request
}

func WriteProtoResponse(w http.ResponseWriter, response *Rsp) {
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)
	data, _ := proto.Marshal(response)
	_, err := w.Write(data)
	if err != nil {
		LogInst().Error(err)
	}
}

func ReadJsonRequest(r *http.Request, val any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, val)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonRequest(w http.ResponseWriter, val any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(val)
	w.Write(data)
}

func rtaIDToSlotPos(rtaID int64) uint16 {
	idStr := fmt.Sprintf("%d", rtaID)
	return uint16(hashKey64(idStr, BitSlotSize))
}

func bToMb(b uint64) uint64 {
	return b >> 20
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func hashKey64(key string, size int) uint64 {
	hasher := fnv.New64()
	hasher.Write([]byte(key))
	return hasher.Sum64() % uint64(size)
}
