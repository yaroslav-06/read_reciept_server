package server

import (
	"fmt"
	"log"
	"net/http"
)

func create_new_img(sys *IdSys, name string, ip string, w http.ResponseWriter) {
	id := sys.gr.GetNewId()
	sys.id_to_email[id] = NewTrackedEmail(name, ip)

	log.Printf("created image pushed id: %s\n", id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, id)
}
