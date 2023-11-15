package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCreateRequest_UnmarshalError(t *testing.T) {
	assert := assert.New(t)
	httpReq, _ := http.NewRequest("POST", "testingURL", bytes.NewBuffer([]byte("bad json")))
	createReq, err := parseCreateRequest(httpReq)

	assert.NotNil(err)
	assert.Equal("unable to unmarshal request", err.Error())
	assert.Equal(createReq.ServiceName, "")
	assert.Equal(createReq.ServiceID, "")
	assert.Equal(createReq.ID, "")
	assert.Equal(createReq.Model, "")
	assert.Equal(createReq.Type, "")
	assert.Equal(createReq.TS, "")
	assert.Equal(createReq.Severity, "")
	assert.Equal(createReq.TeamSlack, "")
}

func TestParseCreateRequest_InvalidTS(t *testing.T) {
	assert := assert.New(t)
	expectedRequest := CreateAlertRequest{
		ServiceID:   "serviceID_1",
		ServiceName: "serviceName_1",
		Alert: Alert{
			ID:        "alert1",
			Model:     "model1",
			Type:      "type1",
			TS:        "this is wrong",
			Severity:  "extreme",
			TeamSlack: "test_slack",
		},
	}
	rawExpectedReq, _ := json.Marshal(expectedRequest)
	httpReq, _ := http.NewRequest("POST", "testingURL", bytes.NewBuffer(rawExpectedReq))
	createReq, err := parseCreateRequest(httpReq)

	assert.NotNil(err)
	assert.Equal(createReq.ServiceName, "")
	assert.Equal(createReq.ServiceID, "")
	assert.Equal(createReq.ID, "")
	assert.Equal(createReq.Model, "")
	assert.Equal(createReq.Type, "")
	assert.Equal(createReq.TS, "")
	assert.Equal(createReq.Severity, "")
	assert.Equal(createReq.TeamSlack, "")
}

func TestParseCreateRequest_ValidRequest(t *testing.T) {
	assert := assert.New(t)
	expectedRequest := CreateAlertRequest{
		ServiceID:   "serviceID_1",
		ServiceName: "serviceName_1",
		Alert: Alert{
			ID:        "alert1",
			Model:     "model1",
			Type:      "type1",
			TS:        "12341234123",
			Severity:  "extreme",
			TeamSlack: "test_slack",
		},
	}
	rawExpectedReq, _ := json.Marshal(expectedRequest)
	httpReq, _ := http.NewRequest("POST", "testingURL", bytes.NewBuffer(rawExpectedReq))
	createReq, err := parseCreateRequest(httpReq)

	assert.Nil(err)
	assert.Equal(expectedRequest.ServiceID, createReq.ServiceID)
	assert.Equal(expectedRequest.ServiceName, createReq.ServiceName)
	assert.Equal(expectedRequest.ID, createReq.ID)
	assert.Equal(expectedRequest.Model, createReq.Model)
	assert.Equal(expectedRequest.Type, createReq.Type)
	assert.Equal(expectedRequest.TS, createReq.TS)
	assert.Equal(expectedRequest.Severity, createReq.Severity)
	assert.Equal(expectedRequest.TeamSlack, createReq.TeamSlack)
}
