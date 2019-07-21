package pizza

import (
	"encoding/json"
	"net/http"
)

type TrackInfos []TrackInfo

type TrackInfo struct {
	StoreID               string      `json:"StoreID"`
	OrderID               string      `json:"OrderID"`
	OrderDescription      string      `json:"OrderDescription"`
	OrderTakeCompleteTime string      `json:"OrderTakeCompleteTime"`
	AdvancedOrderTime     interface{} `json:"AdvancedOrderTime"`
	Actions               struct {
		Track string `json:"Track"`
	} `json:"Actions"`
}

type TrackerStatus struct {
	StoreAsOfTime         string `json:"StoreAsOfTime,omitempty"`
	StoreID               string `json:"StoreID,omitempty"`
	OrderID               string `json:"OrderID,omitempty"`
	PulseOrderGUID        string `json:"PulseOrderGuid,omitempty"`
	Phone                 string `json:"Phone,omitempty"`
	ServiceMethod         string `json:"ServiceMethod,omitempty"`
	OrderDescription      string `json:"OrderDescription,omitempty"`
	OrderTakeCompleteTime string `json:"OrderTakeCompleteTime,omitempty"`
	TakeTimeSecs          int    `json:"TakeTimeSecs,omitempty"`
	OrderSourceCode       string `json:"OrderSourceCode,omitempty"`
	OrderStatus           string `json:"OrderStatus,omitempty"`
	StartTime             string `json:"StartTime,omitempty"`
	MakeTimeSecs          int    `json:"MakeTimeSecs,omitempty"`
	OvenTime              string `json:"OvenTime,omitempty"`
	OvenTimeSecs          int    `json:"OvenTimeSecs,omitempty"`
	RackTime              string `json:"RackTime,omitempty"`
	RackTimeSecs          int    `json:"RackTimeSecs,omitempty"`
	RouteTime             string `json:"RouteTime,omitempty"`
	DriverID              string `json:"DriverID,omitempty"`
	DriverName            string `json:"DriverName,omitempty"`
	OrderDeliveryTimeSecs int    `json:"OrderDeliveryTimeSecs,omitempty"`
	DeliveryTime          string `json:"DeliveryTime,omitempty"`
	OrderKey              string `json:"OrderKey,omitempty"`
	ManagerID             string `json:"ManagerID,omitempty"`
	ManagerName           string `json:"ManagerName,omitempty"`
}

func (c *Client) GetTrackingUrl(phone string) (string, error) {
	url := trackBaseURL + "/v2/orders"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// hardcoded for US, sorry :(
	req.Header.Set("DPZ-Market", "UNITED_STATES")
	req.Header.Set("DPZ-Language", "en")

	q := req.URL.Query()
	q.Add("phonenumber", phone)
	req.URL.RawQuery = q.Encode()
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	trackInfos := &TrackInfos{}
	if err := json.NewDecoder(resp.Body).Decode(trackInfos); err != nil {
		return "", err
	}

	return (*trackInfos)[0].Actions.Track, nil
}

func (c *Client) Track(path string) (*TrackerStatus, error) {
	url := trackBaseURL + "/v2/orders"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// hardcoded for US, sorry :(
	req.Header.Set("DPZ-Market", "UNITED_STATES")
	req.Header.Set("DPZ-Language", "en")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	trackerStatus := &TrackerStatus{}
	if err := json.NewDecoder(resp.Body).Decode(trackerStatus); err != nil {
		return nil, err
	}

	return trackerStatus, nil
}
