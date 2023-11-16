package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (ah *AlertHandler) createAlert(resp http.ResponseWriter, req *http.Request) {
	createReq, err := parseCreateRequest(req)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
	}

	alert, err := convertCreateRequestToAlert(createReq)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
	}
	responseBody := NonGetSuccessResponse{
		AlertID: alert.ID,
	}
	rawResponeBody, _ := json.Marshal(responseBody)

	if alerts, ok := ah.AlertsByService[createReq.ID]; ok {
		// Assume alerts are unique and not updateable. If we get a request for one that already exists
		// ...return an errors letting user know the alert was not created
		for _, existingAlert := range alerts.Alerts {
			if existingAlert.ID == alert.ID {
				responseBody.AlertID = ""
				responseBody.Error = "alert not created because an alert with the same ID already exists"
				rawResponeBody, _ = json.Marshal(responseBody)
				resp.Write(rawResponeBody)
				resp.Header().Set(ContentTypeHeader, ApplicationJson)
				resp.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		alerts.Alerts = append(alerts.Alerts, alert)
	} else {
		ah.AlertsByService[createReq.ID] = AlertsWithService{
			ServiceID:   createReq.ServiceID,
			ServiceName: createReq.ServiceName,
			Alerts: []Alert{
				alert,
			},
		}
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(rawResponeBody)
	resp.Header().Set(ContentTypeHeader, ApplicationJson)
}

func parseCreateRequest(req *http.Request) (CreateAlertRequest, error) {
	body := req.Body
	defer body.Close()
	rawBody, err := io.ReadAll(body)
	if err != nil {
		return CreateAlertRequest{}, err
	}
	fmt.Println(string(rawBody))

	createReq := CreateAlertRequest{}
	err = json.Unmarshal(rawBody, &createReq)
	if err != nil {
		err = errors.New("unable to unmarshal request")
	}
	return createReq, err
}

func convertCreateRequestToAlert(createReq CreateAlertRequest) (Alert, error) {
	alert := Alert{
		ID:        createReq.ID,
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
