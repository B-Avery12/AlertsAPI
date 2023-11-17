package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (ah *AlertHandler) createAlert(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set(ContentTypeHeader, ApplicationJson)
	createReq, err := parseCreateRequest(req)
	if err != nil {
		respBody := NonGetSuccessResponse{
			Error: err.Error(),
		}
		rawRespBody, _ := json.Marshal(respBody)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(rawRespBody)
		return
	}

	incomingAlert, err := convertCreateRequestToAlert(createReq)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		respBody := NonGetSuccessResponse{
			Error: err.Error(),
		}
		rawRespBody, _ := json.Marshal(respBody)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write(rawRespBody)
		return
	}
	responseBody := NonGetSuccessResponse{
		AlertID: incomingAlert.ID,
	}
	rawResponeBody, _ := json.Marshal(responseBody)
	if alertsWithService, ok := ah.AlertsByService[createReq.ServiceID]; ok {
		// Assume alerts are unique and not updateable. If we get a request for one that already exists
		// ...return an errors letting user know the alert was not created
		for _, existingAlert := range alertsWithService.Alerts {
			if existingAlert.ID == incomingAlert.ID {
				responseBody.AlertID = ""
				responseBody.Error = "alert not created because an alert with the same ID already exists"
				rawResponeBody, _ = json.Marshal(responseBody)
				resp.WriteHeader(http.StatusBadRequest)
				resp.Write(rawResponeBody)
				return
			}
		}
		alertsWithService.Alerts = append(alertsWithService.Alerts, incomingAlert)
		ah.AlertsByService[createReq.ServiceID] = alertsWithService
	} else {
		ah.AlertsByService[createReq.ServiceID] = AlertsWithService{
			ServiceID:   createReq.ServiceID,
			ServiceName: createReq.ServiceName,
			Alerts: []Alert{
				incomingAlert,
			},
		}
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(rawResponeBody)
}

func parseCreateRequest(req *http.Request) (CreateAlertRequest, error) {
	body := req.Body
	defer body.Close()
	rawBody, err := io.ReadAll(body)
	if err != nil {
		return CreateAlertRequest{}, err
	}

	createReq := CreateAlertRequest{}
	err = json.Unmarshal(rawBody, &createReq)
	if err != nil {
		err = errors.New("unable to unmarshal request")
	}
	return createReq, err
}

func convertCreateRequestToAlert(createReq CreateAlertRequest) (Alert, error) {
	alert := Alert{
		ID:        createReq.AlertID,
		Model:     createReq.Model,
		Type:      createReq.Type,
		Severity:  createReq.Severity,
		TeamSlack: createReq.TeamSlack,
	}

	unixTime, err := strconv.Atoi(createReq.TS)
	if err != nil {
		return Alert{}, err
	}

	alert.TS = time.Unix(int64(unixTime), 0)
	return alert, nil
}
