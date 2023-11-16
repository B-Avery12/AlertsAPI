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

func (ah *AlertHandler) createAlert(req *http.Request) error {
	createReq, err := parseCreateRequest(req)
	if err != nil {
		return err
	}

	alert, err := convertCreateRequestToAlert(createReq)
	if err != nil {
		return err
	}

	serviceKey := ServiceKey{ServiceID: createReq.ServiceID, ServiceName: createReq.ServiceName}
	ah.AlertsByService[serviceKey] = append(ah.AlertsByService[serviceKey], alert)
	return nil
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
