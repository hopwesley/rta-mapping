package main

import (
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

func writeResponse(w http.ResponseWriter, response *common.Response) {
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
	writeResponse(w, response)
}

func rtaHintV2(w http.ResponseWriter, r *http.Request) {
	var request = readProtoRequest(w, r)
	var response = workV2.QueryRtaMap(request)
	writeResponse(w, response)
}

func NewHttpService() *Service {
	var s = &Service{}
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: common.LogInst(), NoColor: true}))
	r.Use(middleware.Recoverer)
	//r.MethodFunc(http.MethodPost, "/rta_api", rtaHint)

	r.Route("/rta_api", func(r chi.Router) {
		r.Post("/", rtaHint)
		r.Post("/V1", rtaHint)
		r.Post("/V2", rtaHintV2)
	})
	s.router = r
	return s
}
