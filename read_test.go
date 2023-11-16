package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractQueryParams_ValidParams(t *testing.T) {
	httpReq, _ := http.NewRequest("POST", "testingURL", nil)
	queries := httpReq.URL.Query()
	serviceID, startTS, endTS := "my_test_service_id", 1695643160, 1695644360
	queries.Add(ServiceIDParam, serviceID)
	queries.Add(StartTSParam, "1695643160")
	queries.Add(EndTSParam, "1695644360")
	httpReq.URL.RawQuery = queries.Encode()

	params, err := extractQueryParams(httpReq)

	assert.Nil(t, err)
	assert.Equal(t, serviceID, params.ServiceID)
	assert.Equal(t, time.Unix(int64(startTS), 0), params.StartTS)
	assert.Equal(t, time.Unix(int64(endTS), 0), params.EndTS)
}

func TestExtractQueryParams_InValidStartTS(t *testing.T) {
	httpReq, _ := http.NewRequest("POST", "testingURL", nil)
	queries := httpReq.URL.Query()
	serviceID := "my_test_service_id"
	queries.Add(ServiceIDParam, serviceID)
	queries.Add(StartTSParam, "not a unix timestamp")
	queries.Add(EndTSParam, "1695644360")
	httpReq.URL.RawQuery = queries.Encode()

	params, err := extractQueryParams(httpReq)

	assert.NotNil(t, err)
	assert.Equal(t, params.ServiceID, "")
	assert.True(t, params.StartTS.IsZero())
	assert.True(t, params.EndTS.IsZero())
}

func TestExtractQueryParams_InValidEndTS(t *testing.T) {
	httpReq, _ := http.NewRequest("POST", "testingURL", nil)
	queries := httpReq.URL.Query()
	serviceID := "my_test_service_id"
	queries.Add(ServiceIDParam, serviceID)
	queries.Add(StartTSParam, "1695644360")
	queries.Add(EndTSParam, "not a unix timestamp")
	httpReq.URL.RawQuery = queries.Encode()

	params, err := extractQueryParams(httpReq)

	assert.NotNil(t, err)
	assert.Equal(t, params.ServiceID, "")
	assert.True(t, params.StartTS.IsZero())
	assert.True(t, params.EndTS.IsZero())
}

func TestExtractQueryParams_EndTSBeforStartTS(t *testing.T) {
	httpReq, _ := http.NewRequest("POST", "testingURL", nil)
	queries := httpReq.URL.Query()
	serviceID := "my_test_service_id"
	queries.Add(ServiceIDParam, serviceID)
	queries.Add(StartTSParam, "1695644360")
	queries.Add(EndTSParam, "1695644310")
	httpReq.URL.RawQuery = queries.Encode()

	params, err := extractQueryParams(httpReq)

	assert.NotNil(t, err)
	assert.Equal(t,
		"invalid query params, end timestamp must be after start timestamp", err.Error())
	assert.Equal(t, params.ServiceID, "")
	assert.True(t, params.StartTS.IsZero())
	assert.True(t, params.EndTS.IsZero())
}

func TestFilterAlerts_AllAlertsFallInRange(t *testing.T) {
	alertsByService := make(map[string]AlertsWithService, 0)
	alertsByService["serviceID_test"] = AlertsWithService{
		ServiceID:   "serviceID_test",
		ServiceName: "serviceName_test",
		Alerts: []Alert{
			{
				ID:        "alert_id",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
			{
				ID:        "alert_id_2",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
			{
				ID:        "alert_id_3",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
		},
	}

	params := GetAlertsParams{
		ServiceID: "serviceID_test",
		StartTS:   time.Unix(1695644359, 0),
		EndTS:     time.Unix(1695644375, 0),
	}

	filteredAlerts := filterAlerts(params, alertsByService)

	assert.Equal(t, "serviceID_test", filteredAlerts.ServiceID)
	assert.Equal(t, "serviceName_test", filteredAlerts.ServiceName)
	assert.Equal(t, 3, len(filteredAlerts.Alerts))
}

func TestFilterAlerts_SomeAlertsFallInRange(t *testing.T) {
	alertsByService := make(map[string]AlertsWithService, 0)
	alertsByService["serviceID_test"] = AlertsWithService{
		ServiceID:   "serviceID_test",
		ServiceName: "serviceName_test",
		Alerts: []Alert{
			{
				ID:        "alert_id",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
			{
				ID:        "alert_id_2",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
			{
				ID:        "alert_id_3",
				TeamSlack: "slack",
				TS:        time.Unix(1695644390, 0),
			},
		},
	}

	params := GetAlertsParams{
		ServiceID: "serviceID_test",
		StartTS:   time.Unix(1695644359, 0),
		EndTS:     time.Unix(1695644375, 0),
	}

	filteredAlerts := filterAlerts(params, alertsByService)

	assert.Equal(t, "serviceID_test", filteredAlerts.ServiceID)
	assert.Equal(t, "serviceName_test", filteredAlerts.ServiceName)
	assert.Equal(t, 2, len(filteredAlerts.Alerts))
}

func TestFilterAlerts_NoAlertsTimePeriod(t *testing.T) {
	alertsByService := make(map[string]AlertsWithService, 0)
	alertsByService["serviceID_test"] = AlertsWithService{
		ServiceID:   "serviceID_test",
		ServiceName: "serviceName_test",
		Alerts: []Alert{
			{
				ID:        "alert_id",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
			{
				ID:        "alert_id_2",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
			{
				ID:        "alert_id_3",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
		},
	}

	params := GetAlertsParams{
		ServiceID: "serviceID_test",
		StartTS:   time.Unix(1695644365, 0),
		EndTS:     time.Unix(1695644375, 0),
	}

	filteredAlerts := filterAlerts(params, alertsByService)

	assert.Equal(t, "serviceID_test", filteredAlerts.ServiceID)
	assert.Equal(t, "serviceName_test", filteredAlerts.ServiceName)
	assert.Equal(t, 0, len(filteredAlerts.Alerts))
}

func TestFilterAlerts_NoAlertsForServiceID(t *testing.T) {
	alertsByService := make(map[string]AlertsWithService, 0)
	alertsByService["serviceID_test"] = AlertsWithService{
		ServiceID:   "serviceID_test",
		ServiceName: "serviceName_test",
		Alerts: []Alert{
			{
				ID:        "alert_id",
				TeamSlack: "slack",
				TS:        time.Unix(1695644360, 0),
			},
		},
	}
}
