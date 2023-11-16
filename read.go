package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const (
	// 	service_id: The identifier of the service for which alerts are requested.
	// start_ts: The starting timestamp epoch of the time period.
	// end_ts: The ending timestamp epoch of the time period.
	ServiceIDParam = "service_id"
	StartTSParam   = "start_ts"
	EndTSParam     = "end_ts"
)

func (ah *AlertHandler) getAlerts(resp http.ResponseWriter, req *http.Request) {
	params, err := extractQueryParams(req)
	if err != nil {
		respBody := NonGetSuccessResponse{
			Error: "invalid query parameters",
		}
		rawRespBody, err := json.Marshal(respBody)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.WriteHeader(http.StatusBadRequest)
		resp.Header().Set(ContentTypeHeader, ApplicationJson)
		resp.Write(rawRespBody)
		return
	}

	alerts := filterAlerts(params, ah.AlertsByService)
	resp.WriteHeader(http.StatusAccepted)
}

func extractQueryParams(req *http.Request) (GetAlertsParams, error) {
	serviceID := req.URL.Query().Get(ServiceIDParam)
	startTS := req.URL.Query().Get(StartTSParam)
	endTS := req.URL.Query().Get(EndTSParam)

	alertParams := GetAlertsParams{
		ServiceID: serviceID,
	}

	startUnixTime, err := strconv.Atoi(startTS)
	if err != nil {
		return GetAlertsParams{}, err
	}
	alertParams.StartTS = time.Unix(int64(startUnixTime), 0)

	endtUnixTime, err := strconv.Atoi(endTS)
	if err != nil {
		return GetAlertsParams{}, err
	}
	alertParams.EndTS = time.Unix(int64(endtUnixTime), 0)

	return alertParams, nil
}

func filterAlerts(params GetAlertsParams, alertsByService map[string]AlertsWithService) []Alert {
	if alerts, ok := alertsByService[params.ServiceID]; ok {
		// This is not correct
		return alerts.Alerts
	}
}
