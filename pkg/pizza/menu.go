package pizza

type Menu struct {
	Meta     MenuMeta             `json:"Misc"`
	Flavors  map[string]ObjectMap `json:"Flavors"`
	Products map[string]Product   `json:"Products"`
	Sizes    map[string]ObjectMap `json:"Sizes"`
	Toppings map[string]ObjectMap `json:"Toppings"`
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
// type Coupons struct{}
// type Sides map[string]Side
