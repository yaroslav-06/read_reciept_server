package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"readrecieptserver/internal/telegram"
)

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
