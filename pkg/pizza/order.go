package pizza

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"errors"
)

type OrderRequest struct {
	Order *Order `json:"Order"`
}

type OrderProduct struct {
	ID                 int                    `json:"ID,omitempty"`
	Code               string                 `json:"Code,omitempty"`
	Qty                int                    `json:"Qty,omitempty"`
	CategoryCode       string                 `json:"CategoryCode,omitempty"`
	FlavorCode         string                 `json:"FlavorCode,omitempty"`
	Status             int                    `json:"Status,omitempty"`
	LikeProductID      int                    `json:"LikeProductID,omitempty"`
	Name               string                 `json:"Name,omitempty"`
	IsNew              bool                   `json:"IsNew,omitempty"`
	NeedsCustomization bool                   `json:"NeedsCustomization,omitempty"`
	AutoRemove         bool                   `json:"AutoRemove,omitempty"`
	Fulfilled          bool                   `json:"Fulfilled,omitempty"`
	Options            map[string]*Option     `json:"Options,omitempty"`
	Tags               map[string]interface{} `json:"Tags,omitempty"`
	Descriptions       []*Description         `json:"descriptions,omitempty"`
}

type Description struct {
	PortionCode string `json:"portionCode"`
	Value       string `json:"value"`
}

type Option map[string]string

type Order struct {
	Address               *Address               `json:"Address,omitempty"`
	Coupons               []*ObjectCode          `json:"Coupons,omitempty"`
	CustomerID            string                 `json:"CustomerID,omitempty"`
	Email                 string                 `json:"Email,omitempty"`
	Extension             string                 `json:"Extension,omitempty"`
	FirstName             string                 `json:"FirstName,omitempty"`
	LastName              string                 `json:"LastName,omitempty"`
	LanguageCode          string                 `json:"LanguageCode,omitempty"`
	OrderChannel          string                 `json:"OrderChannel,omitempty"`
	OrderID               string                 `json:"OrderID,omitempty"`
	OrderMethod           string                 `json:"OrderMethod,omitempty"`
	OrderTaker            *string                `json:"OrderTaker,omitempty"`
	Payments              []*Payment             `json:"Payments,omitempty"`
	Phone                 string                 `json:"Phone,omitempty"`
	Products              []*OrderProduct        `json:"Products,omitempty"`
	Market                string                 `json:"Market,omitempty"`
	Currency              string                 `json:"Currency,omitempty"`
	ServiceMethod         string                 `json:"ServiceMethod,omitempty"`
	SourceOrganizationURI string                 `json:"SourceOrganizationURI,omitempty"`
	StoreID               string                 `json:"StoreID,omitempty"`
	Tags                  map[string]interface{} `json:"Tags,omitempty"`
	Version               string                 `json:"Version,omitempty"`
	NoCombine             bool                   `json:"NoCombine,omitempty"`
	Partners              map[string]interface{} `json:"Partners,omitempty"`
	NewUser               bool                   `json:"NewUser,omitempty"`
	MetaData              map[string]interface{} `json:"metaData,omitempty"`
	Amounts               Amounts                `json:"Amounts,omitempty"`
	BusinessDate          string                 `json:"BusinessDate,omitempty"`
	EstimatedWaitMinutes  string                 `json:"EstimatedWaitMinutes,omitempty"`
	PriceOrderTime        string                 `json:"PriceOrderTime,omitempty"`
	Status                int                    `json:"Status,omitempty"`
	StatusItems           []*ObjectInfo          `json:"StatusItems,omitempty"`
}

type ObjectCode struct {
	Code string `json:"Code"`
}

type Amounts map[string]float64

type Payment struct {
	Type         string  `json:"Type"`
	Amount       float64 `json:"Amount,omitempty"`
	CardType     string  `json:"CardType,omitempty"`
	Number       string  `json:"Number,omitempty"`
	Expiration   string  `json:"Expiration,omitempty"`
	SecurityCode string  `json:"SecurityCode,omitempty"`
	PostalCode   string  `json:"PostalCode,omitempty"`
}

func NewOrder() *Order {
	return &Order{
		Market:                "UNITED_STATES",
		Currency:              "USD",
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

func (order *Order) WithCustomerInfo(firstName, lastName, email string) *Order {
	order.FirstName = firstName
	order.LastName = lastName
	order.Email = email
	return order
}

func (order *Order) WithPhoneNumber(phone string) *Order {
	order.Phone = phone
	return order
}

func (order *Order) WithStoreID(id string) *Order {
	order.StoreID = id
	return order
}

func (order *Order) AddProduct(product *OrderProduct) {
	order.Products = append(order.Products, product)
}

func (order *Order) AddCoupon(code string) {
	if code != "" {
		order.Coupons = append(order.Coupons, &ObjectCode{Code: code})
	}
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

	if len(returnedOrder.Order.Products) == 0 || len(returnedOrder.Order.Products) != len(order.Products) {
		return 0, errors.New("not all products were returned, possibly invalid product submitted")
	}

	if returnedOrder.Order.Status != 1 {
		return 0, buildValidationError(returnedOrder.Order.StatusItems)
	}

	return returnedOrder.Order.Amounts["Customer"], nil
}

func (c *Client) PlaceOrder(order *Order) (*Order, error) {
	request := &OrderRequest{Order: order}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(request); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", orderURL, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", refererURL)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	returnedOrder := &OrderRequest{}
	if err := json.NewDecoder(resp.Body).Decode(returnedOrder); err != nil {
		return nil, err
	}

	return returnedOrder.Order, nil
}

func buildValidationError(statusItems []*ObjectInfo) error {
	var err string
	for _, status := range statusItems {
		if status.Code != "AutoAddedOrderId" {
			err = fmt.Sprintf("%s%q ", err, status.Code)
		}
	}

	return fmt.Errorf("invalid order, status codes: %s", strings.TrimSpace(err))
}
