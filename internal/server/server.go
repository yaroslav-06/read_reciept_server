package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Start(port string) {
	sys := NewIdSys()

	img_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		ip := getIp(r)

		handle_read(sys, r.RequestURI[len("/image/"):], ip)
		send_image(w)
	})

	get_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		ip := getIp(r)

		get_info(sys, r.RequestURI[len("/get/"):], ip, w)
	})

	new_img_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		ip := getIp(r)

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
