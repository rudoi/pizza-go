package pizza

import (
	"encoding/json"
	"net/url"
)

type Store struct {
	StoreID         string      `json:"StoreID"`
	IsDeliveryStore bool        `json:"IsDeliveryStore"`
	IsOpen          bool        `json:"IsOpen"`
	OpenStatus      *OpenStatus `json:"ServiceIsOpen"`
}

type OpenStatus struct {
	Carryout bool `json:"Carryout"`
	Delivery bool `json:"Delivery"`
}

type StoresResponse struct {
	Status           int      `json:"Status"`
	Stores           []*Store `json:"Stores"`
	RequestedAddress *Address `json:"Address"`
}

type Address struct {
	Street       string `json:"Street"`
	StreetNumber string `json:"StreetNumber"`
	StreetName   string `json:"StreetName"`
	UnitType     string `json:"UnitType"`
	UnitNumber   string `json:"UnitNumber"`
	City         string `json:"City"`
	Region       string `json:"Region"`
	PostalCode   string `json:"PostalCode"`
}

func (c *Client) GetNearestStore(address *Address) (*Store, error) {
	url, err := url.Parse(storesURL)
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Add("s", address.Street)
	q.Add("c", address.PostalCode)
	q.Add("type", "Delivery")
	url.RawQuery = q.Encode()

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	storesResponse := &StoresResponse{}
	if err := json.NewDecoder(resp.Body).Decode(storesResponse); err != nil {
		return nil, err
	}

	return storesResponse.Stores[0], nil
}
