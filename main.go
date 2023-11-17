package main

import (
	"net/http"
)

const (
	ApplicationJson   = "application/json"
	ContentTypeHeader = "content-type"
)

func main() {
	alertHandler := AlertHandler{
		make(map[string]AlertsWithService),
	}
	http.HandleFunc("/alerts", alertHandler.handleRequests)

	http.ListenAndServe(":5000", nil)
}

type AlertHandler struct {
	AlertsByService map[string]AlertsWithService
}

func (ah *AlertHandler) handleRequests(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		ah.createAlert(resp, req)
	case http.MethodGet:
		ah.getAlerts(resp, req)
	default:
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
