package main

import (
	"bytes"
	"context"
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
