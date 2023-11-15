package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (ah *AlertHandler) createAlert(req http.Request) error {
	createReq, err := parseCreateRequest(req)
	if err != nil {
		return err
	}

	serviceKey := ServiceKey{ServiceID: createReq.ServiceID, ServiceName: createReq.ServiceName}
	ah.AlertsByService[serviceKey] = append(ah.AlertsByService[serviceKey], createReq.Alert)
	return nil
}

func parseCreateRequest(req http.Request) (CreateAlertRequest, error) {
	body := req.Body
	defer body.Close()
	rawBody, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("error reading req body")
		return CreateAlertRequest{}, err
	}
	fmt.Println(string(rawBody))

	createReq := CreateAlertRequest{}
	err = json.Unmarshal(rawBody, &createReq)
	return createReq, err
}
