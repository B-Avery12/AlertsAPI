package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

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
		ID:          "alert1",
		Model:       "model1",
		Type:        "type1",
		// TS:        "this is wrong",
		Severity:  "extreme",
		TeamSlack: "test_slack",
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
		ServiceID:   "my_test_service_id",
		ServiceName: "my_test_service",
		ID:          "b950482e9911ec7e41f7ca5e5d9a424f",
		Model:       "my_test_model",
		Type:        "anomaly",
		TS:          "1695644160",
		Severity:    "warning",
		TeamSlack:   "slack_ch",
	}

	rawReq := []byte(`{
		"alert_id": "b950482e9911ec7e41f7ca5e5d9a424f",
		"service_id": "my_test_service_id",
		"service_name": "my_test_service",
		"model": "my_test_model",
		"alert_type": "anomaly",
		"alert_ts": "1695644160",
		"severity": "warning",
		"team_slack": "slack_ch"
	   }`)

	httpReq, _ := http.NewRequest("POST", "testingURL", bytes.NewBuffer(rawReq))
	actualReq, err := parseCreateRequest(httpReq)

	assert.Nil(err)
	assert.Equal(expectedRequest.ServiceID, actualReq.ServiceID)
	assert.Equal(expectedRequest.ServiceName, actualReq.ServiceName)
	assert.Equal(expectedRequest.ID, actualReq.ID)
	assert.Equal(expectedRequest.Model, actualReq.Model)
	assert.Equal(expectedRequest.Type, actualReq.Type)
	assert.Equal(expectedRequest.TS, actualReq.TS)
	assert.Equal(expectedRequest.Severity, actualReq.Severity)
	assert.Equal(expectedRequest.TeamSlack, actualReq.TeamSlack)
}

func TestConvertCreateRequest_ValidTS(t *testing.T) {
	assert := assert.New(t)
	createRequest := CreateAlertRequest{
		ServiceID:   "my_test_service_id",
		ServiceName: "my_test_service",
		ID:          "b950482e9911ec7e41f7ca5e5d9a424f",
		Model:       "my_test_model",
		Type:        "anomaly",
		TS:          "1695644160",
		Severity:    "warning",
		TeamSlack:   "slack_ch",
	}
	alert, err := convertCreateRequestToAlert(createRequest)

	assert.Nil(err)
	assert.Equal(createRequest.ID, alert.ID)
	assert.Equal(createRequest.Model, alert.Model)
	assert.Equal(createRequest.Type, alert.Type)
	assert.Equal(time.Time(time.Date(2023, time.September, 25, 6, 16, 0, 0, time.Local)), alert.TS)
	assert.Equal(createRequest.Severity, alert.Severity)
	assert.Equal(createRequest.TeamSlack, alert.TeamSlack)
}

func TestConvertCreateRequest_InvalidTS(t *testing.T) {
	assert := assert.New(t)
	createRequest := CreateAlertRequest{
		ServiceID:   "my_test_service_id",
		ServiceName: "my_test_service",
		ID:          "b950482e9911ec7e41f7ca5e5d9a424f",
		Model:       "my_test_model",
		Type:        "anomaly",
		TS:          "this is invalid",
		Severity:    "warning",
		TeamSlack:   "slack_ch",
	}
	alert, err := convertCreateRequestToAlert(createRequest)

	assert.NotNil(err)
	assert.Equal("", alert.ID)
	assert.Equal("", alert.Model)
	assert.Equal("", alert.Type)
	assert.True(alert.TS.IsZero())
	assert.Equal("", alert.Severity)
	assert.Equal("", alert.TeamSlack)
}
