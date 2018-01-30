package trades

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strings"
)

// HANDLERS
func GetAll(w http.ResponseWriter, r *http.Request) {
	read := &readAllOp{resp: make(chan map[string]Trade)}
	readsAll <- read
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(<-read.resp)
}

func Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	read := &readOneOp{
		key:  strings.ToUpper(params["pair"]),
		resp: make(chan Trade)}
	readsOne <- read
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(<-read.resp)
}
