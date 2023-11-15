package main

type ServiceKey struct {
	ServiceID   string
	ServiceName string
}

type CreateAlertRequest struct {
	Alert
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
}

type Alert struct {
	ID        string `json:"alert_id"`
	Model     string `json:"model"`
	Type      string `json:"alert_type"`
	TS        string `json:"alert_ts"` // Need to use regex to veify this is only digits and no other characters
	Severity  string `json:"severity"`
	TeamSlack string `json:"team_slack"`
}

type ReadAlertsResponse struct {
	Alerts []Alert `json:"alerts"`
}
