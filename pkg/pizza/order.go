package pizza

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

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

func buildValidationError(statusItems []*ObjectInfo) error {
	var err string
	for _, status := range statusItems {
		if status.Code != "AutoAddedOrderId" {
			err = fmt.Sprintf("%s%q ", err, status.Code)
		}
	}

	return fmt.Errorf("invalid order, status codes: %s", strings.TrimSpace(err))
}
