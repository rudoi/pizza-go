package pizza

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
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
}

type Amounts struct{}

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

func (c *Client) ValidateOrder(order *Order) error {
	request := &OrderRequest{Order: order}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(request); err != nil {
		return err
	}

	resp, err := c.Post(validationURL, "application/json", b)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	returnedOrder := &OrderRequest{}
	if err := json.NewDecoder(resp.Body).Decode(returnedOrder); err != nil {
		return err
	}

	if returnedOrder.Order.Status != 1 {
		return errors.Wrapf(errors.New("order invalid, please confirm input"), "returned order: %+v", returnedOrder)
	}

	if len(returnedOrder.Order.Products) == 0 || len(returnedOrder.Order.Products) != len(order.Products) {
		return errors.New("not all products were returned, possibly invalid product submitted")
	}

	return nil
}
