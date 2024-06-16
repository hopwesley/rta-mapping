package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hopwesley/rta-mapping/common"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"testing"
)

func doHttp(url, cTyp string, data []byte) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", cTyp)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status OK, got %v", resp.Status)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
func TestRtaHint(t *testing.T) {
	req := &common.Req{
		Device: &common.Device{
			ImeiMd5:      "15d35cced5fb9fff20fbfe76ac626df6",
			Oaid:         "0hnMDId1xWZ7h8XByOnORe1mu7XlLZFY",
			AndroidIdMd5: "puQAjHRIDGNtEHtxv77eZdm4I6NQJQcD",
		},
		ReqId:  "xxx-xxx-xxx",
		RtaIds: []int64{10003, 10004, 10005},
	}
	rtaTest(t, req)

	req = &common.Req{
		Device: &common.Device{
			ImeiMd5: "299e707b2808f16422884afd66dd3f61",
			Idfa:    "9gDPCuAWkNrvWTl6Z0h5aEVbI6TNZvZQ",
		},
		ReqId:  "xxx-xxx-xxx",
		RtaIds: []int64{10003, 10004, 10005},
	}
	rtaTest(t, req)

	req = &common.Req{
		Device: &common.Device{
			ImeiMd5: "83ddfcbf6eac1f3d9de08d4aa0db54b1",
			Oaid:    "YsVnMHzORfcKTkUtdFykcVtcyP6BWKpe",
		},
		ReqId:  "xxx-xxx-xxx",
		RtaIds: []int64{10002, 10003, 10004},
	}
	rtaTest(t, req)

}
func rtaTest(t *testing.T, req *common.Req) {

	reqData, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	api := "http://localhost:8801" + "/rta_hint"

	respData, err := doHttp(api, "application/x-protobuf", reqData)
	if err != nil {
		t.Fatalf("http failed:%v", err)
	}
	var response common.Rsp
	err = proto.Unmarshal(respData, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	fmt.Println(response.String())
}

func TestRtaUpdate(t *testing.T) {
	var req = &common.RtaUpdateItem{
		RtaID:   10002,
		UserIDs: []int{10000000001, 10000000002, 10000000003},
		IsDel:   false,
	}

	api := "http://localhost:8801" + "/rta_update"
	reqData, _ := json.Marshal(req)
	respData, err := doHttp(api, "application/json", reqData)
	if err != nil {
		t.Fatalf("http failed:%v", err)
	}
	var rsp common.JsonResponse
	err = json.Unmarshal(respData, &rsp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}

func TestIdUpdate(t *testing.T) {

	var request []*common.IDUpdateReq

	var req = &common.IDUpdateReq{
		UserID: 10000000000,
		IMEIMD5: &common.IDOpItem{
			Val:   "15d35cced5fb9fff20fbfe76ac626df6",
			OpTyp: common.IDOpTypAdd,
		},
		OAID: &common.IDOpItem{
			Val:   "0hnMDId1xWZ7h8XByOnORe1mu7XlLZFY",
			OpTyp: common.IDOpTypUpdate,
		},
		IDFA: &common.IDOpItem{
			Val:   "9gDPCuAWkNrvWTl6Z0h5aEVbI6TNZvZQ",
			OpTyp: common.IDOpTypDel,
		},
		AndroidIDMD5: &common.IDOpItem{
			Val:   "puQAjHRIDGNtEHtxv77eZdm4I6NQJQcD",
			OpTyp: common.IDOpTypAdd,
		},
	}
	request = append(request, req)
	req = &common.IDUpdateReq{
		UserID: 10000000001,
		IMEIMD5: &common.IDOpItem{
			Val:   "15d35cced5fb9fff20fbfe76ac626df6",
			OpTyp: common.IDOpTypAdd,
		},
		OAID: &common.IDOpItem{
			Val:   "0hnMDId1xWZ7h8XByOnORe1mu7XlLZFY",
			OpTyp: common.IDOpTypUpdate,
		},
		IDFA: &common.IDOpItem{
			Val:   "9gDPCuAWkNrvWTl6Z0h5aEVbI6TNZvZQ",
			OpTyp: common.IDOpTypDel,
		},
		AndroidIDMD5: &common.IDOpItem{
			Val:   "puQAjHRIDGNtEHtxv77eZdm4I6NQJQcD",
			OpTyp: common.IDOpTypAdd,
		},
	}
	request = append(request, req)
	api := "http://localhost:8801" + "/id_map_update"
	reqData, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	respData, err := doHttp(api, "application/json", reqData)
	if err != nil {
		t.Fatalf("http failed:%v", err)
	}
	var rsp common.JsonResponse
	err = json.Unmarshal(respData, &rsp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}

func TestIdQuery(t *testing.T) {
	var req = &common.JsonRequest{
		//IMEIMD5: "bdcee62535db4000a7b8b5aeec4fc783",
		IDFA: "vQZui7zkOsy7VLf79nmV4xOssOETzCs7",
		OAID: "0hnMDId1xWZ7h8XByOnORe1mu7XlLZFY",
		//AndroidIDMD5: "puQAjHRIDGNtEHtxv77eZdm4I6NQJQcD",
	}

	api := "http://localhost:8801" + "/query_id"
	reqData, _ := json.Marshal(req)
	respData, err := doHttp(api, "application/json", reqData)
	if err != nil {
		t.Fatalf("http failed:%v", err)
	}
	var rsp common.JsonResponse
	err = json.Unmarshal(respData, &rsp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}

func TestRtaQuery(t *testing.T) {
	var req = &common.JsonRequest{
		UserID: 10000000004,
	}

	api := "http://localhost:8801" + "/query_rta"
	reqData, _ := json.Marshal(req)
	respData, err := doHttp(api, "application/json", reqData)
	if err != nil {
		t.Fatalf("http failed:%v", err)
	}
	var rsp common.JsonResponse
	err = json.Unmarshal(respData, &rsp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rsp)
}
