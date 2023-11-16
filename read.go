package main

import (
	"encoding/json"
	"errors"
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

	filteredAlerts := filterAlerts(params, ah.AlertsByService)

	rawResponseBody, err := json.Marshal(filteredAlerts)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(rawResponseBody)
	resp.Header().Set(ContentTypeHeader, ApplicationJson)
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

	if alertParams.EndTS.Before(alertParams.StartTS) {
		return GetAlertsParams{}, errors.New("invalid query params, end timestamp must be after start timestamp")
	}

	return alertParams, nil
}

func filterAlerts(params GetAlertsParams, alertsByService map[string]AlertsWithService) AlertsWithService {
	filteredAlerts := make([]Alert, 0)
	alertsWithService := AlertsWithService{
		ServiceID: params.ServiceID,
	}

	if alerts, ok := alertsByService[params.ServiceID]; ok {
		for _, alert := range alerts.Alerts {
			if (alert.TS.After(params.StartTS) || alert.TS.Equal(params.StartTS)) &&
				(alert.TS.Before(params.EndTS) || alert.TS.Equal(params.EndTS)) {
				filteredAlerts = append(filteredAlerts, alert)
			}
		}
		alertsWithService.ServiceName = alerts.ServiceName
		alertsWithService.Alerts = filteredAlerts
	}
	return alertsWithService
}
