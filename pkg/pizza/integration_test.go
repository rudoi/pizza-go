// +build integration

package pizza

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var client *Client

var _ = BeforeSuite(func() {
	client = &Client{
		Client: http.Client{Timeout: 10 * time.Second},
	}
})

var _ = Describe("pricing flow", func() {
	Context("given an address and a product", func() {
		It("should return a price", func() {

			// big pink building in
			// downtown Portland, OR
			address := &Address{
				Street:     "111 SW 5th Ave",
				PostalCode: "97204",
			}

			product := &OrderProduct{Code: "14SCREEN", Qty: 1}

			By("looking up the nearest store")
			store, err := client.GetNearestStore(address)
			Expect(err).ToNot(HaveOccurred())
			Expect(store).ToNot(BeNil())
			Expect(store.StoreID).To(Equal("7229"))

			By("building an order")
			order := NewOrder().WithAddress(address).WithStoreID(store.StoreID)
			order.AddProduct(product)

			By("pricing the order")
			price, err := client.ValidateOrder(order)
			Expect(err).ToNot(HaveOccurred())
			Expect(price).ToNot(BeZero())
		})
	})
})
