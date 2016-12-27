package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafael-azevedo/outageapi/database"
	"github.com/rafael-azevedo/outageapi/utils"
)

//OutageList GET /outage
func OutageList(w http.ResponseWriter, r *http.Request) {

	var m database.MultiStatus
	err := m.OutageStatus()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err = json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

}

//AssignNode POST /assign
func AssignNode(w http.ResponseWriter, r *http.Request) {

	or := utils.OutageRequest{}
	or, err := utils.CreateRequest(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	or.Action = "assign"
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(or); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go or.Assign()

}

//DeassignNode POST /deassign
func DeassignNode(w http.ResponseWriter, r *http.Request) {
	or := utils.OutageRequest{}
	or, err := utils.CreateRequest(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	or.Action = "deassign"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(or); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go or.Deassign()

}
func GetRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := utils.IsNumeric(vars["id"])
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	or, err := utils.ReadOutageLog(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(or); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
