package main

import (
	"encoding/json"
	"net/http"
)

const (
	ApplicationJson   = "application/json"
	ContentTypeHeader = "content-type"
)

func main() {
	alertHandler := AlertHandler{}
	// NEED TO UPDATE THIS TO USE
	http.HandleFunc("/hello", alertHandler.createAlert)

	http.ListenAndServe(":5000", nil)
}

type AlertHandler struct {
	AlertsByService map[string]AlertsWithService
}

func (ah *AlertHandler) handleRequests(resp http.ResponseWriter, req *http.Request) {
	switch method := req.Method; method {
	case http.MethodPost:
		ah.createAlert(resp, req)
	case http.MethodGet:
		ah.getAlerts(resp, req)
	default:
		respBody := NonGetSuccessResponse{
			AlertID: "",
			Error:   "requested method is not supported",
		}
		rawRespBody, _ := json.Marshal(respBody)
		resp.WriteHeader(http.StatusMethodNotAllowed)
		resp.Write(rawRespBody)
		resp.Header().Set(ContentTypeHeader, ApplicationJson)
	}
}
