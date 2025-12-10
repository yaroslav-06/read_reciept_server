package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func get_info(sys *IdSys, id string, ip string, w http.ResponseWriter) {
	tre, exists := sys.id_to_email[id]
	if !exists {
		log.Println("get id not found")
		return
	}
	rs := make(map[string]int)
	for sip, nm := range tre.reads_from_ip {
		if sip == ip || nm < 1 {
			continue
		}
		rs[sip] = nm
	}
	bt, err := json.Marshal(rs)
	if err != nil {
		log.Printf("marshaling error at get_info: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bt)
}
