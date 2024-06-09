package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hopwesley/rta-mapping/common"
	"github.com/hopwesley/rta-mapping/workV1"
	"github.com/hopwesley/rta-mapping/workV2"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

type Service struct {
	router *chi.Mux
}

func (s *Service) Start() {
	if _sysConfig.UseSSL {
		panic(http.ListenAndServeTLS(":"+_sysConfig.SrvPort, _sysConfig.SSLCertFile, _sysConfig.SSLKeyFile, s.router))

	} else {
		panic(http.ListenAndServe(":"+_sysConfig.SrvPort, s.router))
	}
}

func readProtoRequest(w http.ResponseWriter, r *http.Request) *common.Request {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return nil
	}

	var request = &common.Request{}
	if err := proto.Unmarshal(body, request); err != nil {
		common.LogInst().Error("request is invalid:", err)
		return nil
	}

	return request
}

func writeProtoResponse(w http.ResponseWriter, response *common.Response) {
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.WriteHeader(http.StatusOK)
	data, _ := proto.Marshal(response)
	_, err := w.Write(data)
	if err != nil {
		common.LogInst().Error(err)
	}
}

func rtaHint(w http.ResponseWriter, r *http.Request) {
	var request = readProtoRequest(w, r)
	var response = workV1.QueryRtaMap(request)
	writeProtoResponse(w, response)
}

func rtaHintV2(w http.ResponseWriter, r *http.Request) {
	var request = readProtoRequest(w, r)
	var response = workV2.QueryRtaMap(request)
	writeProtoResponse(w, response)
}

func rtaUpdate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}

	var request = &common.IDUpdateRequest{}
	err = json.Unmarshal(body, request)
	if err != nil {
		http.Error(w, "Invalid update request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response = common.IDMapInst().UpdateIDMap(request)
	data, _ := json.Marshal(response)
	w.Write(data)
}

func ratRelationUpdate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}

	var request = &workV1.RtaUpdateRequest{}
	err = json.Unmarshal(body, request)
	if err != nil {
		http.Error(w, "Invalid update request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response = workV1.IDRatInst().UpdateHintList(request)
	data, _ := json.Marshal(response)
	w.Write(data)
}

func ratRelationUpdateV2(w http.ResponseWriter, r *http.Request) {
}

func NewHttpService() *Service {
	var s = &Service{}
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: common.LogInst(), NoColor: true}))
	r.Use(middleware.Recoverer)

	r.Route("/rta_api", func(r chi.Router) {
		r.Post("/", rtaHint)
		r.Post("/V1", rtaHint)
		r.Post("/V2", rtaHintV2)
	})

	r.MethodFunc(http.MethodPost, "/rta_update", rtaUpdate)

	r.Route("/id_map_update", func(r chi.Router) {
		r.Post("/", ratRelationUpdate)
		r.Post("/V1", ratRelationUpdate)
		r.Post("/V2", ratRelationUpdateV2)
	})

	s.router = r
	return s
}
