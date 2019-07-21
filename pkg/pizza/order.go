package pizza

import (
	"bytes"
	"encoding/json"

	"errors"
)

type OrderRequest struct {
	Order *Order `json:"Order"`
}

type OrderProduct struct {
	ID                 int                    `json:"ID"`
	Code               string                 `json:"Code"`
	Qty                int                    `json:"Qty"`
	CategoryCode       string                 `json:"CategoryCode"`
	FlavorCode         string                 `json:"FlavorCode"`
	Status             int                    `json:"Status"`
	LikeProductID      int                    `json:"LikeProductID"`
	Name               string                 `json:"Name"`
	IsNew              bool                   `json:"IsNew"`
	NeedsCustomization bool                   `json:"NeedsCustomization"`
	AutoRemove         bool                   `json:"AutoRemove"`
	Fulfilled          bool                   `json:"Fulfilled"`
	Options            map[string]*Option     `json:"Options"`
	Tags               map[string]interface{} `json:"Tags"`
	Descriptions       []*Description         `json:"descriptions"`
}

type Description struct {
	PortionCode string `json:"portionCode"`
	Value       string `json:"value"`
}

type Option map[string]string

type Order struct {
	Address               *Address               `json:"Address"`
	Coupons               []interface{}          `json:"Coupons"`
	CustomerID            string                 `json:"CustomerID"`
	Email                 string                 `json:"Email"`
	Extension             string                 `json:"Extension"`
	FirstName             string                 `json:"FirstName"`
	LastName              string                 `json:"LastName"`
	LanguageCode          string                 `json:"LanguageCode"`
	OrderChannel          string                 `json:"OrderChannel"`
	OrderID               string                 `json:"OrderID"`
	OrderMethod           string                 `json:"OrderMethod"`
	OrderTaker            *string                `json:"OrderTaker"`
	Payments              []*Payment             `json:"Payments"`
	Phone                 string                 `json:"Phone"`
	Products              []*OrderProduct        `json:"Products"`
	Market                string                 `json:"Market"`
	Currency              string                 `json:"Currency"`
	ServiceMethod         string                 `json:"ServiceMethod"`
	SourceOrganizationURI string                 `json:"SourceOrganizationURI"`
	StoreID               string                 `json:"StoreID"`
	Tags                  map[string]interface{} `json:"Tags"`
	Version               string                 `json:"Version"`
	NoCombine             bool                   `json:"NoCombine"`
	Partners              map[string]interface{} `json:"Partners"`
	NewUser               bool                   `json:"NewUser"`
	MetaData              map[string]interface{} `json:"metaData"`
	Amounts               Amounts                `json:"Amounts"`
	BusinessDate          string                 `json:"BusinessDate"`
	EstimatedWaitMinutes  string                 `json:"EstimatedWaitMinutes"`
	PriceOrderTime        string                 `json:"PriceOrderTime"`
	Status                int                    `json:"Status"`
	StatusItems           []*ObjectInfo          `json:"StatusItems"`
}

type Amounts map[string]float64

type Payment struct{}

func NewOrder() *Order {
	return &Order{
		LanguageCode:          "en",
		OrderChannel:          "OLO",
		OrderMethod:           "Web",
		ServiceMethod:         "Delivery",
		SourceOrganizationURI: "order.dominos.com",
		Version:               "1.0",
		NoCombine:             true,
		NewUser:               true,
	}
}

func (order *Order) WithAddress(address *Address) *Order {
	order.Address = address
	return order
}

func (order *Order) WithStoreID(id string) *Order {
	order.StoreID = id
	return order
}

func (order *Order) AddProduct(product *OrderProduct) {
	order.Products = append(order.Products, product)
}

// ValidateOrder validates the order and returns its price
func (c *Client) ValidateOrder(order *Order) (float64, error) {
	request := &OrderRequest{Order: order}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(request); err != nil {
		return 0, err
	}

	resp, err := c.Post(pricingURL, "application/json", b)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	returnedOrder := &OrderRequest{}
	if err := json.NewDecoder(resp.Body).Decode(returnedOrder); err != nil {
		return 0, err
	}

	if returnedOrder.Order.Status != 1 {
		return 0, errors.New("order invalid, please confirm input")
	}

	if len(returnedOrder.Order.Products) == 0 || len(returnedOrder.Order.Products) != len(order.Products) {
		return 0, errors.New("not all products were returned, possibly invalid product submitted")
	}

	for _, item := range returnedOrder.Order.StatusItems {
		if item.Code == "BelowMinimumDeliveryAmount" {
			return 0, errors.New("order does not meet minimum delivery price")
		}
	}

	return returnedOrder.Order.Amounts["Customer"], nil
}
