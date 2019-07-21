package pizza

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Menu struct {
	Meta     MenuMeta               `json:"Misc"`
	Flavors  map[string]*ObjectMap  `json:"Flavors"`
	Products map[string]*Product    `json:"Products"`
	Sizes    map[string]*ObjectMap  `json:"Sizes"`
	Toppings map[string]*ObjectMap  `json:"Toppings"`
	Variants map[string]*Variant    `json:"Variants"`
	Coupons  map[string]*ObjectInfo `json:"Coupons"`
}

type MenuMeta struct {
	Status        int    `json:"Status"`
	StoreID       string `json:"StoreID"`
	BusinessDate  string `json:"BusinessDate"`
	StoreAsOfTime string `json:"StoreAsOfTime"`
	LanguageCode  string `json:"LanguageCode"`
	Version       string `json:"Version"`
	ExpiresOn     string `json:"ExpiresOn"`
}

type ObjectMap map[string]ObjectInfo

type ObjectInfo struct {
	Code        string                 `json:"Code"`
	ImageCode   string                 `json:"ImageCode"`
	Description string                 `json:"Description"`
	Local       bool                   `json:"Local"`
	Name        string                 `json:"Name"`
	SortSeq     string                 `json:"SortSeq"`
	Tags        map[string]interface{} `json:"Tags"`
}

type Product struct {
	AvailableToppings string                 `json:"AvailableToppings"`
	AvailableSides    string                 `json:"AvailableSides"`
	Code              string                 `json:"Code"`
	DefaultToppings   string                 `json:"DefaultToppings"`
	DefaultSides      string                 `json:"DefaultSides"`
	Description       string                 `json:"Description"`
	ImageCode         string                 `json:"ImageCode"`
	Local             bool                   `json:"Local"`
	Name              string                 `json:"Name"`
	ProductType       string                 `json:"ProductType"`
	Tags              map[string]interface{} `json:"Tags"`
	Variants          []string               `json:"Variants"`
}

type Variant struct {
	Code                       string                 `json:"Code"`
	FlavorCode                 string                 `json:"FlavorCode"`
	ImageCode                  string                 `json:"ImageCode"`
	Local                      bool                   `json:"Local"`
	Name                       string                 `json:"Name"`
	Price                      string                 `json:"Price"`
	ProductCode                string                 `json:"ProductCode"`
	SizeCode                   string                 `json:"SizeCode"`
	Tags                       map[string]interface{} `json:"Tags"`
	AllowedCookingInstructions string                 `json:"AllowedCookingInstructions"`
	DefaultCookingInstructions string                 `json:"DefaultCookingInstructions"`
	Prepared                   bool                   `json:"Prepared"`
}

// not sure if needed yet
// type Categorization struct{}
// type Sides map[string]Side

func (c *Client) GetStoreMenu(storeID string) (*Menu, error) {
	url, err := url.Parse(fmt.Sprintf(menuURL, storeID))
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Add("lang", "en")
	q.Add("structured", "true")
	url.RawQuery = q.Encode()

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	menu := &Menu{}
	if err := json.NewDecoder(resp.Body).Decode(menu); err != nil {
		return nil, err
	}

	return menu, nil
}

// Rather than include a ton of lookup logic for coupons,
// just look up the ridiculous 50% off coupon that sometimes
// exists.
func (m *Menu) GetFiftyPercentCouponCode() string {
	for _, coupon := range m.Coupons {
		if coupon.ImageCode == "OLO50" {
			return coupon.Code
		}
	}

	return ""
}
