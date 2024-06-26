package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hopwesley/rta-mapping/common"
	"net/http"
)

type Service struct {
	router *chi.Mux
}

func (s *Service) Start() {
	if _sysConfig.UseSSL {
		fmt.Println("https service start success:", _sysConfig.SrvPort)
		panic(http.ListenAndServeTLS(":"+_sysConfig.SrvPort, _sysConfig.SSLCertFile, _sysConfig.SSLKeyFile, s.router))
	} else {
		fmt.Println("http service start success:", _sysConfig.SrvPort)
		panic(http.ListenAndServe(":"+_sysConfig.SrvPort, s.router))
	}
}

func rtaHint(w http.ResponseWriter, r *http.Request) {
	var request = common.ReadProtoRequest(w, r)
	var response = common.CheckIfHinted(request)
	common.WriteProtoResponse(w, response)
}

func rtaUpdate(w http.ResponseWriter, r *http.Request) {
	var req = &common.RtaUpdateItem{}
	err := common.ReadJsonRequest(r, req)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	var res = common.RtaMapInst().UpdateRta(req)
	common.WriteJsonRequest(w, res)
}

func idUpdate(w http.ResponseWriter, r *http.Request) {
	var request []*common.IDUpdateReq
	err := common.ReadJsonRequest(r, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var response = common.IDMapInst().UpdateIDMap(request)
	common.WriteJsonRequest(w, response)
}

func rtaQuery(w http.ResponseWriter, r *http.Request) {
	var request = &common.JsonRequest{}
	err := common.ReadJsonRequest(r, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var response = common.RtaMapInst().QueryRatInfos(request)
	common.WriteJsonRequest(w, response)
}

func idQuery(w http.ResponseWriter, r *http.Request) {
	var request = &common.JsonRequest{}
	err := common.ReadJsonRequest(r, request)
	if err != nil {
		http.Error(w, "Invalid query request", http.StatusBadRequest)
		return
	}
	var response = common.IDMapInst().QueryIDInfos(request)
	common.WriteJsonRequest(w, response)
}

func NewHttpService() *Service {
	var s = &Service{}
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: common.LogInst(), NoColor: true}))
	r.Use(middleware.Recoverer)

	r.MethodFunc(http.MethodPost, "/rta_hint", rtaHint)
	r.MethodFunc(http.MethodPost, "/rta_update", rtaUpdate)

	r.MethodFunc(http.MethodPost, "/id_map_update", idUpdate)

	r.MethodFunc(http.MethodPost, "/query_rta", rtaQuery)
	r.MethodFunc(http.MethodPost, "/query_id", idQuery)

	s.router = r
	return s
}
