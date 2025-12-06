package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"readrecieptserver/internal/telegram"
	uniqueid "readrecieptserver/internal/unique_id"
	"strings"
)

type IdSys struct {
	gr          *uniqueid.Generator
	id_to_email map[string]*TrackedEmail
}

type TrackedEmail struct {
	name          string
	reads_from_ip map[string]int
}

func NewTrackedEmail(name string, ip string) *TrackedEmail {
	rs := &TrackedEmail{}
	rs.name = name
	rs.reads_from_ip = make(map[string]int)
	rs.reads_from_ip[ip] = -1 // to ignore sender ip in future
	return rs
}

func NewIdSys() *IdSys {
	sys := &IdSys{}

	sys.gr = uniqueid.NewGenerator()
	sys.id_to_email = make(map[string]*TrackedEmail)

	return sys
}

func StartServer(port string) {
	sys := NewIdSys()

	img_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		ip := strings.Split(r.RemoteAddr, ":")[0] //get remadd and removes port

		handle_read(sys, r.RequestURI[len("/image/"):], ip)
		send_image(w)
	})

	get_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		ip := strings.Split(r.RemoteAddr, ":")[0] //get remadd and removes port

		get_info(sys, r.RequestURI[len("/get/"):], ip, w)
	})

	new_img_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		ip := strings.Split(r.RemoteAddr, ":")[0] //get remadd and removes port

		bd, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("couldn't read the body")
			return
		}
		create_new_img(sys, string(bd), ip, w)
	})

	http.Handle("/image/", img_handler)
	http.Handle("/get/", get_handler)
	http.Handle("/new_img", new_img_handler)

	fmt.Printf("Server started at port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

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

func create_new_img(sys *IdSys, name string, ip string, w http.ResponseWriter) {
	id := sys.gr.GetNewId()
	sys.id_to_email[id] = NewTrackedEmail(name, ip)

	log.Printf("created image pushed id: %s\n", id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, id)
}

func send_image(w http.ResponseWriter) {
	buf, err := os.ReadFile("assets/pixel.png")
	if err != nil {
		fmt.Printf("couldn't load image: %s", err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(buf)
}

func handle_read(sys *IdSys, id string, ip string) {
	log.Printf("read with id: %s\n", id)
	tre, exists := sys.id_to_email[id]
	if !exists {
		log.Println("couldn't find the id")
		return
	}
	if tre.reads_from_ip[ip] != -1 {
		tre.reads_from_ip[ip] += 1
	}
	telegram.SendMsg("@parolk06", fmt.Sprintf("read the message: %s", tre.name))
}
