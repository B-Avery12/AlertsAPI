package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	ServiceIDParam = "service_id"
	StartTSParam   = "start_ts"
	EndTSParam     = "end_ts"
)

func (ah *AlertHandler) getAlerts(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set(ContentTypeHeader, ApplicationJson)
	params, err := extractQueryParams(req)
	if err != nil {
		respBody := NonGetSuccessResponse{
			Error: fmt.Sprintf("invalid query parameters: %v", err),
		}
		rawRespBody, _ := json.Marshal(respBody)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(rawRespBody)
		return
	}

	filteredAlerts := filterAlerts(params, ah.AlertsByService)

	rawResponseBody, err := json.Marshal(filteredAlerts)
	if err != nil {
		errRespBody := NonGetSuccessResponse{
			Error: "internal error",
		}
		errRawRespBody, _ := json.Marshal(errRespBody)
		resp.Write(errRawRespBody)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(rawResponseBody)
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
		return GetAlertsParams{}, errors.New("end timestamp must be after start timestamp")
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
