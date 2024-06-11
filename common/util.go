package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"math/big"
	"net/http"
)

var (
	Version   string
	Commit    string
	BuildTime string
)

const (
	HitSuccess = iota
	HitFailed
)
const (
	BidTypeOk = iota
)

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
	hash := md5.Sum([]byte(idStr))
	hashInt := new(big.Int)
	hashInt.SetBytes(hash[:])
	modulus := big.NewInt(BitSlotSize)
	remainder := new(big.Int).Mod(hashInt, modulus)
	return uint16(remainder.Int64())
}
