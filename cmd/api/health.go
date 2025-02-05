package main

import (
	"net/http"
)

func (api *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status": "ok",
		"env":    api.config.env,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusOK, err.Error())
	}
}
