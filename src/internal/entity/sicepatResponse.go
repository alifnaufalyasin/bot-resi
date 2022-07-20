package entity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SicepatRes struct {
	Sicepat Sicepat `json:"sicepat"`
}

type Sicepat struct {
	Status Status `json:"status"`
	Result Result `json:"result"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Result struct {
	WaybillNumber      string         `json:"waybill_number"`
	Service            string         `json:"service"`
	ReceiverAddress    string         `json:"receiver_address"`
	ReceiverName       string         `json:"receiver_name"`
	PODReceiver        string         `json:"POD_receiver"`
	TotalPrice         int            `json:"totalprice"`
	RealPrice          int            `json:"realprice"`
	TrackHistory       []TrackHistory `json:"track_history"`
	SendImage          string         `json:"pop_img_path"`
	PersonSendImage    string         `json:"pop_sigesit_img_path"`
	ReceiveImage       string         `json:"pod_img_path"`
	PersonReceiveImage string         `json:"pod_sigesit_img_path"`
}

type TrackHistory struct {
	Date         string `json:"date_time"`
	Status       string `json:"status"`
	City         string `json:"city"`
	ReceiverName string `json:"receiver_name"`
}

var (
	client = http.Client{
		Timeout: time.Second * 10,
	}
)

func GetResiSicepatHistory(uri, resi string) (SicepatRes, string, error) {
	getResponse := SicepatRes{}
	rawBody := ""

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", uri, resi), nil)
	if err != nil {
		return getResponse, rawBody, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("referer", "https://www.sicepat.com/")
	req.Header.Add("origin", "https://www.sicepat.com")

	resp, err := client.Do(req)
	if err != nil {
		return getResponse, rawBody, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return getResponse, rawBody, err
	}

	rawBody = string(body)

	err = json.Unmarshal(body, &getResponse)
	if err != nil {
		return getResponse, rawBody, err
	}
	return getResponse, rawBody, nil
}
