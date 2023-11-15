package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (ah *AlertHandler) createAlert(req *http.Request) error {
	createReq, err := parseCreateRequest(req)
	if err != nil {
		return err
	}

	serviceKey := ServiceKey{ServiceID: createReq.ServiceID, ServiceName: createReq.ServiceName}
	ah.AlertsByService[serviceKey] = append(ah.AlertsByService[serviceKey], createReq.Alert)
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
