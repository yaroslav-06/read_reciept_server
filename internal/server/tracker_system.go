package server

import (
	"net/http"
	uniqueid "readrecieptserver/internal/unique_id"
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

func getIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
