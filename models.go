package main

import "time"

type CreateAlertRequest struct {
	AlertID     string `json:"alert_id"`
	Model       string `json:"model"`
	Type        string `json:"alert_type"`
	TS          string `json:"alert_ts"`
	Severity    string `json:"severity"`
	TeamSlack   string `json:"team_slack"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
}

type Alert struct {
	ID        string    `json:"alert_id"`
	Model     string    `json:"model"`
	Type      string    `json:"alert_type"`
	TS        time.Time `json:"alert_ts"`
	Severity  string    `json:"severity"`
	TeamSlack string    `json:"team_slack"`
}

type AlertsWithService struct {
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	Alerts      []Alert `json:"alerts"`
}

type NonGetSuccessResponse struct {
	AlertID string `json:"alert_id"`
	Error   string `json:"error"`
}

type GetAlertsParams struct {
	ServiceID string
	StartTS   time.Time
	EndTS     time.Time
}
