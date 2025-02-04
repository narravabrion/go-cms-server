package main

import "net/http"

func (api *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("server is ok!"))
}
