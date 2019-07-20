package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rudoi/pizza-go/pkg/pizza"
)

func main() {
	address := &pizza.Address{
		Street:     "111 SW 5th Ave",
		PostalCode: "97204",
	}

	client := pizza.Client{
		Client: http.Client{Timeout: 10 * time.Second},
	}

	store, err := client.GetNearestStore(address)
	if err != nil {
		panic(err)
	}

	fmt.Println(store.StoreID)

	order := pizza.NewOrder().WithAddress(address).WithStoreID(store.StoreID)
	order.AddProduct(&pizza.OrderProduct{Code: "14SCREEN", Qty: 1})

	if err := client.ValidateOrder(order); err != nil {
		panic(err)
	} else {
		fmt.Println("Order validated.")
	}
}
